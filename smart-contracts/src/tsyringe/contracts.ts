import { inject, injectable, instanceCachingFactory, registry, singleton } from "tsyringe"
import type { Contract } from "ethers"
import { BigNumber, ContractFactory } from "ethers"
import { HardhatRuntimeEnvironment } from "hardhat/types"
import { EthereumAddress, NotNativeCurrencyAddress } from "../ethereumAddress"
import { HardhatRuntimeEnvironmentToken, NetworkDescriptorToken } from "./injectionTokens"
import { SifchainAccounts, SifchainAccountsPromise } from "./sifchainAccounts"
import {
  BridgeBank,
  BridgeBank__factory,
  BridgeRegistry,
  BridgeRegistry__factory,
  BridgeToken,
  BridgeToken__factory,
  CosmosBridge__factory,
} from "../../build"
import "@openzeppelin/hardhat-upgrades"

import web3 from "web3"
export const MINTER_ROLE = web3.utils.soliditySha3("MINTER_ROLE")
const ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000"

@singleton()
export class SifchainContractFactories {
  bridgeBank: Promise<BridgeBank__factory>
  cosmosBridge: Promise<CosmosBridge__factory>
  bridgeRegistry: Promise<BridgeRegistry__factory>
  bridgeToken: Promise<BridgeToken__factory>

  constructor(@inject(HardhatRuntimeEnvironmentToken) hre: HardhatRuntimeEnvironment) {
    this.bridgeBank = hre.ethers
      .getContractFactory("BridgeBank")
      .then((x: ContractFactory) => x as BridgeBank__factory)
    this.cosmosBridge = hre.ethers
      .getContractFactory("CosmosBridge")
      .then((x: ContractFactory) => x as CosmosBridge__factory)
    this.bridgeRegistry = hre.ethers
      .getContractFactory("BridgeRegistry")
      .then((x: ContractFactory) => x as BridgeRegistry__factory)
    this.bridgeToken = hre.ethers
      .getContractFactory("BridgeToken")
      .then((x: ContractFactory) => x as BridgeToken__factory)
  }
}

export class CosmosBridgeArguments {
  constructor(
    readonly operator: EthereumAddress,
    readonly consensusThreshold: number,
    readonly initialValidators: Array<EthereumAddress>,
    readonly initialPowers: Array<number>,
    readonly networkDescriptor: number
  ) {}

  asArray() {
    return [
      this.operator.address,
      this.consensusThreshold,
      this.initialValidators.map((x) => x.address),
      this.initialPowers,
      this.networkDescriptor,
    ]
  }
}

export class CosmosBridgeArgumentsPromise {
  constructor(readonly cosmosBridgeArguments: Promise<CosmosBridgeArguments>) {}
}

@singleton()
export class CosmosBridgeProxy {
  contract: Promise<Contract>

  constructor(
    @inject(HardhatRuntimeEnvironmentToken) hardhatRuntimeEnvironment: HardhatRuntimeEnvironment,
    sifchainContractFactories: SifchainContractFactories,
    cosmosBridgeArgumentsPromise: CosmosBridgeArgumentsPromise
  ) {
    this.contract = sifchainContractFactories.cosmosBridge.then(async (cosmosBridgeFactory) => {
      const args = await cosmosBridgeArgumentsPromise.cosmosBridgeArguments
      const cosmosBridgeProxy = await hardhatRuntimeEnvironment.upgrades.deployProxy(
        cosmosBridgeFactory,
        args.asArray(),
        { initializer: "initialize(address,uint256,address[],uint256[],int32)" }
      )
      await cosmosBridgeProxy.deployed()
      return cosmosBridgeProxy
    })
  }
}

export function defaultCosmosBridgeArguments(
  sifchainAccounts: SifchainAccounts,
  power: number = 100,
  networkDescriptor: number = 9999
): CosmosBridgeArguments {
  const powers = sifchainAccounts.validatatorAccounts.map((_) => power)
  const threshold = powers.reduce((acc, x) => acc + x)
  return new CosmosBridgeArguments(
    new NotNativeCurrencyAddress(sifchainAccounts.operatorAccount.address),
    threshold,
    sifchainAccounts.validatatorAccounts.map((x) => new NotNativeCurrencyAddress(x.address)),
    powers,
    networkDescriptor
  )
}

@registry([
  {
    token: CosmosBridgeArgumentsPromise,
    useFactory: instanceCachingFactory<CosmosBridgeArgumentsPromise>((c) => {
      const accountsPromise = c.resolve(SifchainAccountsPromise)
      return new CosmosBridgeArgumentsPromise(
        accountsPromise.accounts.then((accts) => {
          return defaultCosmosBridgeArguments(accts)
        })
      )
    }),
  },
  {
    token: NetworkDescriptorToken,
    useValue: 1,
  },
])
@injectable()
export class BridgeBankArguments {
  constructor(
    private readonly cosmosBridgeProxy: CosmosBridgeProxy,
    private readonly sifchainAccountsPromise: SifchainAccountsPromise,
    @inject(NetworkDescriptorToken) private readonly networkDescriptor: number
  ) {}

  async asArray() {
    const cosmosBridge = await this.cosmosBridgeProxy.contract
    const accts = await this.sifchainAccountsPromise.accounts
    const result = [
      accts.operatorAccount.address,
      cosmosBridge.address,
      accts.ownerAccount.address,
      accts.pauserAccount.address,
      this.networkDescriptor,
    ]
    return result
  }
}

@singleton()
export class BridgeBankProxy {
  contract: Promise<BridgeBank>

  constructor(
    @inject(HardhatRuntimeEnvironmentToken) h: HardhatRuntimeEnvironment,
    private sifchainContractFactories: SifchainContractFactories,
    private bridgeBankArguments: BridgeBankArguments
  ) {
    this.contract = sifchainContractFactories.bridgeBank.then(async (bridgeBankFactory) => {
      const bridgeBankArguments = await this.bridgeBankArguments.asArray()
      const bridgeBankProxy = (await h.upgrades.deployProxy(
        bridgeBankFactory,
        bridgeBankArguments,
        {
          initializer: "initialize(address,address,address,address,int32)",
          unsafeAllow: ["delegatecall"],
        }
      )) as BridgeBank
      await bridgeBankProxy.deployed()
      return bridgeBankProxy
    })
  }
}

@singleton()
export class BridgeRegistryProxy {
  contract: Promise<BridgeRegistry>

  constructor(
    @inject(HardhatRuntimeEnvironmentToken) h: HardhatRuntimeEnvironment,
    private sifchainContractFactories: SifchainContractFactories,
    private cosmosBridgeProxy: CosmosBridgeProxy,
    private bridgeBankProxy: BridgeBankProxy
  ) {
    this.contract = sifchainContractFactories.bridgeRegistry.then(async (bridgeRegistryFactory) => {
      const bridgeRegistryProxy = await h.upgrades.deployProxy(bridgeRegistryFactory, [
        (await cosmosBridgeProxy.contract).address,
        (await bridgeBankProxy.contract).address,
      ])
      await bridgeRegistryProxy.deployed()
      return bridgeRegistryProxy as BridgeRegistry
    })
  }
}

/**
 * Deploys a BridgeToken named erowan
 */
@singleton()
export class RowanContract {
  readonly contract: Promise<BridgeToken>

  constructor(private sifchainContractFactories: SifchainContractFactories) {
    this.contract = sifchainContractFactories.bridgeToken.then(async (bridgeToken) => {
      return (await (bridgeToken as BridgeToken__factory).deploy(
        "erowan",
        "erowan",
        18,
        "rowan"
      )) as BridgeToken
    })
  }
}

@singleton()
export class BridgeTokenSetup {
  readonly complete: Promise<boolean>

  private async build(
    rowan: RowanContract,
    bridgeBankProxy: BridgeBankProxy,
    sifchainAccounts: SifchainAccountsPromise
  ) {
    const erowan = await rowan.contract
    const owner = (await sifchainAccounts.accounts).ownerAccount
    const bridgebank = (await bridgeBankProxy.contract).connect(owner)
    await bridgebank.addExistingBridgeToken(erowan.address)
    await erowan.grantRole(String(MINTER_ROLE), bridgebank.address)
    await erowan.grantRole(String(MINTER_ROLE), owner.address)
    await erowan.grantRole(ADMIN_ROLE, bridgebank.address)
    await erowan.approve(bridgebank.address, "10000000000000000000")
    const accounts = await sifchainAccounts.accounts
    const muchRowan = BigNumber.from(100000000).mul(BigNumber.from(10).pow(18))
    await erowan.connect(owner).mint(accounts.operatorAccount.address, muchRowan)
    return true
  }

  constructor(
    rowan: RowanContract,
    bridgeBankProxy: BridgeBankProxy,
    sifchainAccounts: SifchainAccountsPromise
  ) {
    this.complete = this.build(rowan, bridgeBankProxy, sifchainAccounts)
  }
}

/**
 * TODO: BridgeBank does not function unless this runs. This should be part of
 * BridgeBankProxy. This is a temporary setup because we ran into circular
 * dependency issue when we tried to
 */
@singleton()
export class BridgeBankSetup {
  readonly complete: Promise<boolean>

  private async build(
    bridgebankProxy: BridgeBankProxy,
    cosmosBridgeProxy: CosmosBridgeProxy,
    sifchainAccounts: SifchainAccountsPromise
  ): Promise<boolean> {
    const cosmosBridge = await cosmosBridgeProxy.contract
    const operator = (await sifchainAccounts.accounts).operatorAccount
    const bridgebank = await bridgebankProxy.contract

    await cosmosBridge.connect(operator).setBridgeBank(bridgebank.address)
    return true
  }

  constructor(
    bridgebankProxy: BridgeBankProxy,
    cosmosBridgeProxy: CosmosBridgeProxy,
    sifchainAccounts: SifchainAccountsPromise
  ) {
    this.complete = this.build(bridgebankProxy, cosmosBridgeProxy, sifchainAccounts)
  }
}
