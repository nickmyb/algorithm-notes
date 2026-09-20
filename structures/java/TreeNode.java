/**
 * 二叉树节点。
 *
 * <p>字段和构造器照搬 LeetCode Java 编辑器里给的那段注释（Definition for a binary tree node），
 * 一个字不改。这样题解在本地和在提交框里编译结果一致，可以原样来回复制。
 *
 * <p>建树、遍历这类只有本地测试才需要的辅助方法放在 {@link TreeNodes}，不挂在这里，
 * 挂上去就破坏了上面那条性质。
 *
 * <p>javatest.sh 编译每个题目目录时会把 structures/java 下的文件一并加进源文件列表，
 * 所以题解里直接用 TreeNode 即可。<b>题目目录里不要再定义自己的 TreeNode</b>，会撞类名。
 */
public class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;

    TreeNode() {}

    TreeNode(int val) {
        this.val = val;
    }

    TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}
