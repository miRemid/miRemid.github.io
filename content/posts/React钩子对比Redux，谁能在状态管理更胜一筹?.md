---
title: "React钩子对比Redux，谁能在状态管理更胜一筹?"
date: 2021-06-21T14:47:25+08:00
draft: true
toc: false
images:
tags: 
  - untagged
---

> 本文翻译自[Nicola Grujicic](https://www.framelessgrid.com/react-hooks-vs-redux-for-state-management-in-2021/)

## React状态管理：Hook和Redux，两者有何不同？
世界上很多包括我在内的开发者都有一个疑惑，在听到人们使用Hook代替Redux时，为什么要将一个多级别工作非常完美的组件用一个新的还未完成的东西替代

> 许多人在疑惑useContext和useReducer是否能代替Redux

今天，我将会尝试去回答这个问题。我将会解释React的钩子(useContext,useReducer)和Redux的不同之处以及使用场景

> 两者最大的查表在于如何管理应用的全局状态

Redux是当开发者需要创建一个大型的复杂的应用时管理全局状态的最佳实践，它提供了一个相当于存储中心的组件用于保管需要在整个应用中使用的状态信息。并通过一系列规则保证状态只有在指向性的数据流中被更新

![Redux](https://images.ctfassets.net/yytn7c23rcp1/5IgvC7d2l2CXbkxHm2Syol/b1494f62db8cdc974f85086a8c12a9a1/redux_diagram.png)

Reack钩子则是用另一种全新的方式在组件的生命周期中来管理状态信息并且并不依赖于组件。它在16.8版本中首次提出，旨在通过组件分享逻辑来降低组件之间的复杂性

两者最大的不同就是，Redux创建了一个包裹在整个应用之外的全局状态容器，这个容器被称之为`store`并且通过`useReducer`创建一个独立的组件用于和应用组件进行合作

另外一点是，React钩子使用useContext和useReducer联合去和组件状态管理进行合作，这已经是另一个层面的合作关系。它让通过useReducer创建的状态容器和它的`dispatch`函数能够传递给从最顶端组件以下的任意一个组件之中。你可以在最顶端组件中定义它们使其变为全局状态

如果我们将我们创建的所有state放置在最外层的组件，这体验感觉就和使用Redux没有什么区别

> 所以是Context(上下文)可能代替Redux(但也不一定)

对于Redux，使用钩子的好处有

并不需要持续关注Redux第三方依赖的更新，避免更新所导致的Bug和问题，可以让你的应用大小变得更小也因此会有更快的速度来处理状态变化。可以让代码更为清晰便于理解，并且通过钩子来创建组件可以用更少的代码获取更快的响应速度

尽管看上去使用useReducer和Redux没什么区别，但它并不是Redux。useReducer函数和其reducer紧密联系在一起，使得能够使用它的dispatch函数。但我们仅仅向reducer传递dispatch的action对象而已

__你可以认为Redux是一个全局状态总线，它将承载任意一个事件(action)并且基于action提供的数据和状态进行处理__

## 总结
Redux和其他React状态管理解决方案之所以被提出是因为在React组件中管理一个全局状态是一个非常愚蠢的想法。如果这么做会导致一系列问题，如多源数据和使用虚拟模型代替原始模型展示等等

如果钩子能够让状态管理更清晰，我认为在小型应用中可以尝试不使用Redux。但对于大型的复杂的应用程序，全局状态管理仍然不能被忽略