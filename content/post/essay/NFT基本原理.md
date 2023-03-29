---
title: "NFT基本原理"
date: 2021-06-14T19:18:03+08:00
draft: true
toc: false
images:
tags: 
  - 杂谈
  - 区块链
  - 翻译
---

> 本文节选翻译自[Github Ethereum EIPs Eip-721#rationale](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md#rationale)

现在有很多旨在追踪可区分资产的基于以太坊智能合约的提议。例如Decentraland现有或计划中的的NFT、CryptoPunks中的eponymouse和游戏中如DMarket或EnjinCoin的物品系统。在未来，可以用于追踪现实生活中的资产例如房产（如ubiquity和Propy等公司的想法）。当这些物品在分类账单中没有集中在一起时显得至关重要，每个资产必须单独并能够独立追踪其所有权信息。

## 为什么选择NFT这个词
NFT这个词让几乎每个人都能够满意的接受，并且广泛适用于可区分数字资产的领域。我们认识到契约在对于一个本标准的某些应用时又确切的描述（特别是物理性质）。

除了NFT这个词还考虑了别的单词，例如：Alternatives considered, Distinguishable asset, Title, Token, Asset, Equity, Ticket

## NFT的标识
每一个NFT都由一个使用ERC-721智能合约生成的独一无二的`uint256`ID。这个验证数字在合约中将不会发生改变。这个配对信息（合约地址、uint256 ID）将会成为全球中独一无二且完全合格的以太坊链上的特定资产证明。尽管有些ETC-721智能合约可能发现从0开始依次递增的NFT ID非常方便，但是调用者并不能假设那个ID数字有对他们具有任何特定的模式，并把这些ID成为黑匣子。要注意的是，NFTs可能会变成无效Token或被摧毁。请参照支持的枚举界面中的枚举函数。

使用uint256作为ID让广泛的应用得以运用，因为UUID和sha3哈希散列值可以直接转换为uint256。

## 转移机制
ETC-721定义了一个安全的转移函数`safeTransferForm`（不使用`bytes`参数进行重载）和一个不安全函数`transferFrom`。转移可能会在下述角色中进行初始化：
- NFT的所有者
- 接收NFT的地址
- 一个对于当前NFT拥有者可信任的操作员

另外，可信任的操作员可能会为NFT设置一个可接受的地址。这为钱包、代理人和实际应用提供一系列有力的工具来使用这些大量NFT。

转移和接受函数文档只适用于交易必须抛出时的指定条件下。你自己的复现可能会在其他情况下同样被抛出。这让众多复现产生有趣的结果
- Disallow transfers if the contract is paused — prior art, CryptoKitties deployed contract, line 611
- Blacklist certain address from receiving NFTs — prior art, CryptoKitties deployed contract, lines 565, 566
- Disallow unsafe transfers — transferFrom throws unless _to equals msg.sender or countOf(_to) is non-zero or was non-zero previously (because such cases are safe)
- Charge a fee to both parties of a transaction — require payment when calling approve with a non-zero _approved if it was previously the zero address, refund payment if calling approve with the zero address if it was previously a non-zero address, require payment when calling any transfer function, require transfer parameter _to to equal msg.sender, require transfer parameter _to to be the approved address for the NFT
- Read only NFT registry — always throw from unsafeTransfer, transferFrom, approve and setApprovalForAll

在ETC-223、ERC-677、ERC-827和OpenZeppelin在`SafeERC20.sol`中对NFT的最佳实现中都定义了失败的交易将会被抛弃。ETC-20定义了一个成为`allowance`的特性，这将会导致当交易被唤起时，但在另一个状态下发生改变时发生错误，在OpenZeppelin的#438个讨论中有被提及。在ETC-721中并没有`allowance`，因为每个NFT都是独一无二的，数量为0或1.因此，我们可以吸收ETC-20'初始设计的优点，而且不会产生后续的问题。

制造NFT（铸币）和销毁NFT（燃烧）并不包含在规范中。在你可以通过其他方法来实现这些功能。当你负责创建或销毁NFT时，请参照`event`文档。

我们在思考`onERC721Received`函数中的`operator`参数是否必要。我们可以想象，在所有的情况下，如果在交易中操作员的地位十分重要，操作员可以将token转给自己然后转回，然后他们将会成为交易中的所有者地址。这似乎让我们认为操作员是token的所有者（转移到自己账户的行为是多余的）。当操作员发送token时，它会按照操作员自己的协议行事，而不是代表token持有人行事的操作员。这也是为什么操作员和token的前任收件人拥有的重要意义。

*正在考虑中的替代方案：只允许ERC-20协议的两步验证交易，要求转移函数永远不会抛出，要求所有函数返回一个能够代表操作是否成功的布尔值。*

## ETC-165接口
我们选择了可靠的ETC-165接口来暴露所有ETC-721智能合约支持的接口

一个未来的EIP可能会创造一个全球性的河源登记册接口。我们强烈支持例如例如EIP这种接口，它允许你的ERC-721复现去实现`ETC721Enumerable`,`ETC721Metadata`或其他接口通过委派到单独的合约。

## 费用和复杂性
该规范考虑了管理少数大量NFT的实现。如果你的应用能够迭代，请在你的代码中避免使用for/while循环（参考CryptoKitties赏金问题#4）。这种循环可能让你的合约无法进行扩展，并且在无限制的情况下，合约产生的gas成本将会随时间的推移而上升。

我们在测试网络中部署了一个XXXXERC721合约并追踪了340282366920938463463374607431768211456 不同行为。这足以将每隔IPv6地址分配给以太坊账户所有者，或追踪数量占地球大小一半的纳米机器人。

*考虑的替代方案：如果要使用for循环，删除资产枚举功能，从枚举函数中返回稳定的数组类型。*

## 隐私
钱包、经纪人、拍卖行都在动机上都有强烈的证明NFT所有者的需求。

想想使用NFT不可枚举的使用情况可能会分有意思，例如属性的所有权或部分私有注册表。但是，隐私不可能被实现，因为攻击者可以很轻松通过调用`ownerOf`对所有可能的`tokenId`进行信息检索。

## 元数据的选择
我们已经要求了在metadata扩展中需要`name`和`symbl`函数。每一个我们看过的Token规划和蓝图中（ETC-20、ETC-223、ETC-677、ETC-777、ETC-827）都包含这些函数。

我们提醒开发者，如果不想使用这个机制，即便是空的字符串也是合法的名称或符号。除此之外，我们提醒大家，任何只能合约都可以使用与合约相同的名称或符号。客户如何去确定哪些ETC-721智能合约是规范的超出了本标准的范围。

有一种机制提供了在NFTs和URIs之间的关联。我们希望在许多实现中可以利用这一点来为每个NFT提供元数据。图像尺寸推荐取自Instagram，它们可能在图片的可用性方面有很多了解。URI是可变的（不定时更改），我们可以认为NFT代表了一幢房子的所有权，在这种情况下，关于房屋的元数据可以自然地发生改变。

元数据将使用字符串的形式返回。目前，这仅是从Web3的调用中获取，而不是从其他合约中。这一点是可以接受的，因为我们还没有考虑在一个区块内申请查询此类信息的情况。

*考虑的替代方案：在区块链上保存每一个资产的元信息（非常昂贵），使用URL模板来查询元信息部分（URL模板并不适用所有URL协议，尤其是P2P），多地址网络*

## 向下兼容性
我们已经从ETC-20规范中接受了`balanceOf`,`totalSupply`,`name`和`symbol`这些定义。当然在实现中也可以包含用于支持ETC-20自己的协议的返回uint8的`decimals`函数。

截止2018年，有很多NFT的实现例子
- CryptoKitties -- Compatible with an earlier version of this standard.
- CryptoPunks -- Partially ERC-20 compatible, but not easily generalizable because it includes auction functionality directly in the contract and uses function names that explicitly refer to the assets as "punks".
- Auctionhouse Asset Interface -- The author needed a generic interface for the Auctionhouse ÐApp (currently ice-boxed). His "Asset" contract is very simple, but is missing ERC-20 compatibility, approve() functionality, and metadata. This effort is referenced in the discussion for EIP-173.