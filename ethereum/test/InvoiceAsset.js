// const assert = require("assert");
const Web3 = require("web3");
const web3 = new Web3();
const BigNumber = require("bignumber.js");
const ethers = require("ethers");

contract("DV", (accounts) => {
  let owner;

  let DVContract = artifacts.require("DV");
  let ruc20 = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

  let mapper = {
    N: 5,
    NT: [4, 3],
    E: 5,
    P: 7,
    I: 9,
    AV: [1, 5],
  };

  contract("#dv", () => {
    before(async () => {
      owner = accounts[0];
      contract = await DVContract.new();
      await contract.seed();
      assert.equal(contract !== null, true);
    });

    describe("when verifying N identities", () => {
      it("should verify 827400125 to match 91", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...[0, 0],
          ...[2, 7, 4, 0, 0, 1, 2, 5],
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        // console.log("\n [PRINT BIGNUMBER] ");
        assert.equal(joinedResp, "91");
      });

      it("should verify 8274301253 to match 31", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...[0, 0],
          ...[2, 7, 4, 3, 0],
          ...[1, 2, 5, 3],
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "31");
      });

      it("should verify 87132230 to match 11", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...[0, 0],
          ...[7, 1, 3], // 4
          ...[2, 2, 3], // 6
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "11");
      });

      it("should verify 8747704 to match 67", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...[0, 0],
          ...[7, 4, 7, 0, 0], // 4
          ...[7, 0, 4], // 6
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "67");
      });

      it("should verify 8NT0010024 to match 33", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...mapper.NT,
          ...[0, 0, 1, 0, 0, 0, 2, 4],
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "33");
      });

      it("should verify 8NT001123456 to match 76", async () => {
        const id = [
          mapper.N,
          ...[0, 8],
          ...mapper.NT,
          ...[0, 0, 1, 1, 2, 3, 4, 5, 6],
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "76");
      });

      it("should verify PE0010019 to match 60", async () => {
        const id = [
          mapper.N,
          ...[0, 0],
          ...[mapper.P, mapper.E],
          ...[0, 0, 1, 0, 0, 0, 1, 9],
        ];

        let ruc21 = [...ruc20.slice(0, 20 - id.length), ...id, 0];
        const resp = await contract.calc(ruc21.map((i) => new BigNumber(i)));
        const joinedResp = resp.toString().replace(",", "");
        assert.equal(joinedResp, "60");
      });

      it("should withdraw", async () => {
        const resp = await contract.withdraw(owner);
        console.log(resp);
      });
    });
  });
});
