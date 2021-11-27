import Crypto

pub contract AnconVerifier {
    pub resource interface IQueryRoot {
        pub fun queryRootCalculation(
            leafOpUint: [UInt8],
            prefix: [UInt8],
            existenceProofInnerOp: [[[UInt8]]],
            existenceProofInnerOpHash: UInt8,
            existenceProofKey: [UInt8],
            existenceProofValue: [UInt8]
         ): [UInt8];
    }

    pub resource QueryRoot {
        pub fun queryRootCalculation(
            leafOpUint: [UInt8],
            prefix: [UInt8],
            existenceProofInnerOp: [[[UInt8]]],
            existenceProofInnerOpHash: UInt8,
            existenceProofKey: [UInt8],
            existenceProofValue: [UInt8]
         ): [UInt8] {
            return self.queryRootCalculation(
                leafOpUint: leafOpUint,
                prefix: prefix,
                existenceProofInnerOp: existenceProofInnerOp,
                existenceProofInnerOpHash: existenceProofInnerOpHash,
                existenceProofKey: existenceProofKey,
                existenceProofValue: existenceProofValue,
            );
        }
    }

    init() {
        // Publish the "QueryRoot" resource, so that any client can use it.
        self.account.save(<- create QueryRoot(), to: /storage/AnconQueryRoot);
        self.account.link<&{IQueryRoot}>(/public/AnconQueryRoot, target: /storage/AnconQueryRoot);
    }

    // Data structures and helper functions
    pub enum HashOp: UInt8 {
        pub case NO_HASH
        pub case SHA256
        pub case SHA512
        pub case KECCAK
        pub case RIPEMD160
        pub case BITCOIN
    }
    pub enum LengthOp: UInt8 {
        pub case NO_PREFIX
        pub case VAR_PROTO
        pub case VAR_RLP
        pub case FIXED32_BIG
        pub case FIXED32_LITTLE
        pub case FIXED64_BIG
        pub case FIXED64_LITTLE
        pub case REQUIRE_32_BYTES
        pub case  REQUIRE_64_BYTES
    }

   pub struct ExistenceProof {
      pub var  valid: Bool
      pub var  key: [UInt8]
      pub var  value: [UInt8] 
      pub var  leaf: LeafOp
      pub var  path:[InnerOp] 

      init(
       valid: Bool,
       key: [UInt8],
       value: [UInt8] ,
       leaf: LeafOp,
       path:[InnerOp]           
      ){
            self.valid = valid
            self.key = key
            self.value = value
            self.leaf = leaf
            self.path = path
      }
    }

   pub struct NonExistenceProof {
      pub var    valid: Bool
      pub var   key:[UInt8]
      pub var   left:ExistenceProof
      pub var   right:ExistenceProof

      init(
                   valid: Bool,
        key:[UInt8],
        left:ExistenceProof,
        right:ExistenceProof,
      ) {
          self.valid = valid
          self.key = key
          self.left = left
          self.right = right
      }
    }

   pub struct LeafOp {
      pub var    valid: Bool
      pub var   hash:HashOp
      pub var   prehash_key:HashOp
      pub var   prehash_value:HashOp
      pub var   len:LengthOp
      pub var  prefix:  [UInt8]

      init(
          valid: Bool,
         hash:HashOp,
         prehash_key:HashOp,
         prehash_value:HashOp,
         len:LengthOp,
        prefix:  [UInt8]

      ) {
          self.valid=valid
          self.hash=hash
          self.prehash_key=prehash_key
          self.prehash_value=prehash_value
          self.len=len
          self.prefix=prefix
      }
    }

    

   pub struct InnerOp {
      pub var    valid: Bool
     pub var   hash:HashOp
      pub var  prefix:  [UInt8]
      pub var  suffix:  [UInt8]
          init(
          valid: Bool,
         hash:HashOp,
        prefix:  [UInt8],
        suffix:  [UInt8]         

      ) {
          self.valid=valid
          self.hash=hash
          self.suffix=suffix
          self.prefix=prefix
      }
    }



   pub struct ProofSpec {
      pub var leafSpec: LeafOp 
      pub var innerSpec: InnerSpec 
      pub var maxDepth: Int 
      pub var minDepth: Int 
          init(
          leafSpec:  LeafOp,
         innerSpec:InnerSpec,
      ) {
          self.leafSpec=leafSpec
          self.innerSpec=innerSpec
          self.maxDepth=0
          self.minDepth=0
      }
    }



   pub struct InnerSpec {
      pub var childOrder: [Int] 
      pub var childSize: Int 
      pub var minPrefixLength: Int 
      pub var maxPrefixLength: Int 
      pub var emptyChild: [UInt8] 
      pub var hash: HashOp 
          init(
 childOrder: [Int] ,
       childSize: Int ,
       minPrefixLength: Int ,
       maxPrefixLength: Int ,
       emptyChild: [UInt8] ,
       hash: HashOp ,

      ) {
          self.childOrder=childOrder
          self.childSize=childSize
          self.minPrefixLength=minPrefixLength
          self.maxPrefixLength=maxPrefixLength
          self.emptyChild=emptyChild    
          self.hash=hash
      }
    }



   pub fun getIavlSpec()    : ProofSpec  {
       

    let  leafSpec = LeafOp(
         valid:   true,
            hash: HashOp.SHA256,
            prehash_key: HashOp.NO_HASH,
           prehash_value: HashOp.SHA256,
         len:   LengthOp.VAR_PROTO,
prefix: [0 as UInt8]
        );

    let  innerSpec = InnerSpec(
        childOrder:    [0 as Int,1 as Int],
           childSize: 33,
           minPrefixLength: 4,
           maxPrefixLength: 12,
           emptyChild: [],
         hash:   HashOp.SHA256
        );
       return   ProofSpec(leafSpec: leafSpec, innerSpec: innerSpec
        )

    }

    pub enum Ordering: UInt8{
        pub case LT
        pub case EQ
        pub case GT
    }

   pub fun checkAgainstSpecLeafOp(op: LeafOp, spec: ProofSpec)
    {
        assert(op.hash == spec.leafSpec.hash,message:  "Unexpected HashOp");
        assert(
            op.prehash_key == spec.leafSpec.prehash_key,
            message: "Unexpected PrehashKey"
        );
        assert(
            op.prehash_value == spec.leafSpec.prehash_value,
            message: "Unexpected PrehashKey"
        );
        assert(op.len == spec.leafSpec.len,message:  "UnexpecteleafSpec LengthOp");
        let leafSpecPrefix = String.encodeHex(spec.leafSpec.prefix)
        let opPrefix = String.encodeHex(op.prefix)

        assert(
            spec.leafSpec.prefix.contains(op.prefix.removeFirst()),
            message: "LeafOpLib: wrong prefix"
        );
    }

   pub fun applyValueLeafOp(
        op:LeafOp,
        key:[UInt8],
        value:[UInt8],
    )    : [UInt8] {
        assert(key.length > 0, message: "Leaf op needs key");
        assert(value.length > 0,message: "Leaf op needs value");

        let keyLeafData = self.prepareLeafData(hashOp:op.prehash_key, lengthOp :op.len, data: key);
        let valueLeafData = self.prepareLeafData(hashOp:op.prehash_value, lengthOp :op.len, data: value);

        let data = op.prefix.concat(keyLeafData).concat(valueLeafData);
        return self.doHash(hashOp:op.hash, data: data);
    }

   pub fun prepareLeafData(
        hashOp:HashOp,
        lengthOp:LengthOp,
        data:[UInt8]  ,
    )    : [UInt8] {
        return self.doLength(lengthOp: lengthOp, data :self.doHashOrNoop(hashOp: hashOp, data: data));
    }

   pub fun doHashOrNoop(
       hashOp:HashOp , 
       data:[UInt8]   ,
       )    : [UInt8]
    {
        if (hashOp == HashOp.NO_HASH) {
            return data;
        }
        return self.doHash(hashOp: hashOp, data: data);
    }

   pub fun checkAgainstSpecInnerOp(op: InnerOp, spec: ProofSpec)
    {
        assert(op.hash == spec.leafSpec.hash,message:  "Unexpected HashOp");

        assert(
            !spec.leafSpec.prefix.contains(op.prefix[0]),
            message: "LeafOpLib: wrong prefix"
        );
        assert(
            op.prefix.length >= spec.innerSpec.minPrefixLength,
            message: "InnerOp prefix too short"
        );

        let maxLeftChildLen = (spec.innerSpec.childOrder.length - 1) * spec.innerSpec.childSize;
        assert(
            op.prefix.length <=
                (spec.innerSpec.maxPrefixLength) + maxLeftChildLen,
            message: "InnerOp prefix too short"
        );
    }

   pub fun applyValueInnerOp(op: InnerOp, child:[UInt8]): [UInt8]
    {
        assert(child.length > 0, message: "Inner op needs child value");
        return self.doHash(
                hashOp: op.hash,
                data: op.prefix.concat(child).concat(op.suffix)
            );
    }

   pub fun doHash(hashOp: HashOp, data: [UInt8]):[UInt8] 
    {
        // The only available hashing function available right now is SHA256.
        return HashAlgorithm.SHA2_256.hash(data);
    }

   pub fun doLength(lengthOp: LengthOp , data: [UInt8]): [UInt8]
    {
        if (lengthOp == LengthOp.NO_PREFIX) {
            return data;
        }
        if (lengthOp == LengthOp.VAR_PROTO) {
            return self.encodeVarintProto(n:  data.length).concat(data);
        }
        if (lengthOp == LengthOp.REQUIRE_32_BYTES) {
            assert(data.length == 32, message: "Expected 32 [UInt8]");
            return data;
        }
        if (lengthOp == LengthOp.REQUIRE_64_BYTES) {
            assert(data.length == 64, message: "Expected 64 [UInt8]");
            return data;
        }
        panic("Unsupported lengthop");
    }

   pub fun encodeVarintProto(n: Int)    : [UInt8]   {
        // Inspired by: https://github.com/confio/ics23/blob/17c3466679ba7b65e1207406951e50577b4dfb05/js/src/ops.ts#L148
        let enc: [UInt8] = [];
        var l = n;
        while (l >= 128) {
            let b = UInt8((l % 128) + 128);
            enc.append(b);
            l = l / 128;
        }
        enc.append(UInt8(l));
        return enc;
    }

   pub fun verify(
        proof: ExistenceProof,
        spec: ProofSpec,
        root: [UInt8],
        key: [UInt8],
        value: [UInt8],
    )    {
        self.checkAgainstSpec(proof: proof, spec: spec);
        assert(
            String.encodeHex(proof.key) == String.encodeHex(key),
            message: "Provided key doesn't match proof"
        )
        assert(
            String.encodeHex(proof.value) == String.encodeHex(value),
            message: "Provided value doesn't match proof"
        )
        assert(
            String.encodeHex(self.calculate(p: proof)) == String.encodeHex(root),
            message: "Calculcated root doesn't match provided root"
        )
    }

   pub fun checkAgainstSpec(proof: ExistenceProof, spec: ProofSpec)            
    {
        self.checkAgainstSpecLeafOp(op: proof.leaf, spec: spec)

        assert(
            spec.minDepth == 0 || proof.path.length >= (spec.minDepth),
            message: "InnerOps depth too short"
        )
        assert(
            spec.maxDepth == 0 || proof.path.length >= (spec.maxDepth),
            message: "InnerOps depth too short"
        )

        for path in proof.path {
            self.checkAgainstSpecInnerOp(op: path, spec: spec)
        }
    }

    // Calculate determines the root hash that matches the given proof.
    // You must validate the result is what you have in a header.
    // Returns error if the calculations cannot be performed.
   pub fun calculate(p: ExistenceProof) : [UInt8]
    {
        // leaf step takes the key and value as input
        var res = self.applyValueLeafOp(op: p.leaf, key: p.key, value: p.value)

        // the rest just take the output of the last step (reducing it)
        for path in p.path {
            res = self.applyValueInnerOp(op: path, child: res)
        }
        return res
    }

    //Separate arguments to not convert from uint 256 to byte
    pub fun convertProof(
        key: [UInt8],
        value: [UInt8],
        _prefix: [UInt8],
        _leafOpUint: [UInt8],
        _innerOp: [[[UInt8]]],
        existenceProofInnerOpHash: UInt8
    ): ExistenceProof {
        assert(
            _leafOpUint.length >= 4,
            message: "Not enough values packed inside leafOpUInt"
        );
        let leafOp = LeafOp(
            valid: true,
            hash: HashOp(rawValue: _leafOpUint[0])!,
            prehash_key: HashOp(rawValue: _leafOpUint[1])!,
            prehash_value: HashOp(rawValue: _leafOpUint[2])!,
            len: LengthOp(rawValue: _leafOpUint[3])!,
            prefix: _prefix
        );

        let innerOpArr: [InnerOp] = [];

        var i = 0;
        while (i < _innerOp.length) {
            var temp: [[UInt8]] = _innerOp[i];
            innerOpArr.append(InnerOp(
                valid: true,
                hash: HashOp(rawValue: existenceProofInnerOpHash)!,
                prefix: temp[0],
                suffix: temp[1]
            ));
            i = i + 1;
        }

        return ExistenceProof(
            valid: true,
            key: key,
            value: value,
            leaf: leafOp,
            path: innerOpArr
        );
    }

    pub fun queryRootCalculation(
        leafOpUint: [UInt8],
        prefix: [UInt8],
        existenceProofInnerOp: [[[UInt8]]],
        existenceProofInnerOpHash: UInt8,
        existenceProofKey: [UInt8],
        existenceProofValue: [UInt8]
    ): [UInt8] {
        let proof: ExistenceProof = self.convertProof(
            key: existenceProofKey,
            value: existenceProofValue,
            _prefix: prefix,
            _leafOpUint: leafOpUint,
            _innerOp: existenceProofInnerOp,
            existenceProofInnerOpHash: existenceProofInnerOpHash
        );
        return self.calculate(p: proof);
    }

    // claimed ics23 proofs
    // key = prefix + cid eg  ancon+cid
    // value = sha256(dagcbor)
    pub fun changeOwnerWithProof(
        leafOpUint: [UInt8],
        prefix: [UInt8],
        existenceProofInnerOp: [[[UInt8]]],
        existenceProofInnerOpHash: UInt8,
        existenceProofKey: [UInt8],
        existenceProofValue: [UInt8],
        root: [UInt8],
        key: [UInt8],
        value: [UInt8]
    ): Bool {
        let proof = self.convertProof(
            key: existenceProofKey,
            value: existenceProofValue,
            _prefix: prefix,
            _leafOpUint: leafOpUint,
            _innerOp: existenceProofInnerOp,
            existenceProofInnerOpHash: existenceProofInnerOpHash
        );

        self.verify(
            proof: proof,
            spec: self.getIavlSpec(),
            root: root,
            key: key,
            value: value
        );

        return true;
    }

}
