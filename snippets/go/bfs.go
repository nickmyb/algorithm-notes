package snippets

// BFS 骨架。adj(v) 表示 v 的邻居；marked 和 edgeTo 由调用方准备。
// 若不需要还原路径，可以删掉 edgeTo；多源 BFS 则先将所有起点入队并标记。
//
// q := []int{s}
// marked[s] = true
// for len(q) > 0 {
//     v := q[0]
//     q = q[1:]
//     for _, w := range adj(v) {
//         if !marked[w] {
//             q = append(q, w)
//             marked[w] = true
//             edgeTo[w] = v
//         }
//     }
// }
