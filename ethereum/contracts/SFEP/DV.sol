pragma solidity ^0.8.7;
pragma experimental ABIEncoderV2;

import "../../node_modules/@openzeppelin/contracts/utils/Address.sol";

contract DV {
    using Address for address payable;

    event Withdrawn(address indexed payee, uint256 weiAmount);
    event LogDV(uint256[] dv);

    uint256 fee = 0.002 * 1e18;
    address owner;
    mapping(uint256 => mapping(uint256 => uint256)) private chars;
    mapping(uint256 => mapping(uint256 => uint256)) private corresp;

    constructor() public {
        owner = msg.sender;
    }

    function withdraw(address payable payee) public {
        require(msg.sender == owner, "INVALID_USER");
        uint256 b = address(this).balance;
        payee.sendValue(address(this).balance);

        emit Withdrawn(payee, b);
    }

    function setFee(uint256 _fee) public {
        require(msg.sender == owner, "INVALID_USER");
        fee = _fee;
    }

    function getFee() public returns (uint256) {
        return fee;
    }

    function findCategoryAndReplace(uint256[] memory ruc21)
        internal
        view
        returns (bytes32[] memory, uint256[] memory)
    {
        bytes32[] memory c;
        if (getCharsAt(6, ruc21) == 5 || getCharsAt(7, ruc21) == 5) {
            if (getCharsAt(7, ruc21) == 5) {
                ruc21[7] = 5;
            }

            if (getCharsAt(6, ruc21) == 5) {
                ruc21[6] = 5;
            }
            //     c.push(keccak256('N'));
        }

        if (getCharsAt(10, ruc21) == 4 && getCharsAt(11, ruc21) == 3) {
            /// c.push(keccak256('NT'));
            ruc21[10] = 4;
            ruc21[11] = 3;
        }

        if (getCharsAt(10, ruc21) == 5) {
            ruc21[10] = 5;
            //   c.push(keccak256('E'));
        }

        return (c, ruc21);
    }

    function getCharsAt(uint256 pos, uint256[] memory ruc21)
        public
        view
        returns (uint256)
    {
        return chars[pos][ruc21[pos]];
    }

    function seed() public returns (bool) {
        require(msg.sender == owner, "INVALID_USER");

        chars[6][90] = 5; // N
        chars[7][90] = 5; // N

        // NT
        chars[9][90] = 4; // N
        chars[10][91] = 3; // T

        // E
        chars[10][92] = 5; // E

        // ruc18 = E o N
        // skip

        // PE
        chars[9][93] = 7;
        chars[10][92] = 5;

        // PI
        chars[9][93] = 7;
        chars[10][94] = 9;

        // AV
        chars[9][95] = 1;
        chars[10][96] = 5;

        corresp[0][0] = 100;
        corresp[1][0] = 1;
        corresp[1][1] = 2;
        corresp[1][2] = 3;
        corresp[1][3] = 4;
        corresp[1][4] = 5;
        corresp[1][5] = 6;
        corresp[1][6] = 7;
        corresp[1][7] = 8;
        corresp[1][8] = 9;
        corresp[1][9] = 1;
        corresp[2][0] = 2;
        corresp[2][1] = 3;
        corresp[2][2] = 5;
        corresp[2][3] = 7;
        corresp[2][4] = 8;
        corresp[2][5] = 9;
        corresp[2][6] = 2;
        corresp[2][7] = 3;
        corresp[2][8] = 4;
        corresp[2][9] = 5;
        corresp[3][0] = 6;
        corresp[3][1] = 7;
        corresp[3][2] = 8;
        corresp[3][3] = 9;
        corresp[3][4] = 1;
        corresp[3][5] = 2;
        corresp[3][6] = 3;
        corresp[3][7] = 4;
        corresp[3][8] = 5;
        corresp[3][9] = 6;
        corresp[4][0] = 7;
        corresp[4][1] = 8;
        corresp[4][2] = 9;
        corresp[4][3] = 1;
        corresp[4][4] = 2;
        corresp[4][5] = 3;
        corresp[4][6] = 4;
        corresp[4][7] = 5;
        corresp[4][8] = 6;
        corresp[4][9] = 7;
    }

    function isJurArrVal(uint256[] memory ruc21) private view returns (bool) {
        // pos 6 y 7
        uint256 pos = corresp[ruc21[5]][ruc21[6]];
        if (pos > 0) {
            return true;
        } else {
            return false;
        }
    }

    function concat(bytes memory a, bytes memory b)
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodePacked(a, b);
    }

    function calc(uint256[] memory ruc21)
        public
        view
        returns (uint256[] memory)
    {
        // require(msg.value >= fee, "MUST SEND FEE BEFORE USE");

        (bytes32[] memory cat, uint256[] memory ruc21) = findCategoryAndReplace(
            ruc21
        );
        uint256[] memory dv = new uint256[](2);
        // if (isJurArrVal(ruc21)) {
        //     //
        // } else {
        // calcula dv
        (ruc21, dv) = calcDV(ruc21, dv, uint256(19), uint256(0));
        (ruc21, dv) = calcDV(ruc21, dv, uint256(20), uint256(1));

        // emit LogDV(dv);
        return dv;
    }

    function calcDV(
        uint256[] memory ruc,
        uint256[] memory dv,
        uint256 len,
        uint256 cycle
    ) private view returns (uint256[] memory, uint256[] memory) {
        uint256 suma = 0;
        uint256 j = 2;

        uint256 i = len;
        uint256 sw = 0;
        while (i > 0) {
            suma = suma + j * ruc[i];
            j = j + 1;
            i = i - 1;
        }
        uint256 _dv = 0;
        if (suma > 0) {
            uint256 x = suma % 11;

            if (x == 0 || x == 1) {
                _dv = 0;
            } else {
                _dv = uint256(11 - x);
            }
        }

        if (cycle == 0) {
            ruc[20] = _dv;
            dv[0] = _dv;
        }

        if (cycle == 1) {
            dv[1] = _dv;
        }

        return (ruc, dv);
    }
}
