package snippets

// DFS 骨架。adj(v) 表示 v 的邻居；marked 和 edgeTo 由调用方准备。
// 若不需要还原路径，可以删掉 edgeTo。
//
// var dfs func(int)
// dfs = func(v int) {
//     marked[v] = true
//     for _, w := range adj(v) {
//         if !marked[w] {
//             edgeTo[w] = v
//             dfs(w)
//         }
//     }
// }
// dfs(s)
