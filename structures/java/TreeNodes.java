import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;

/**
 * 二叉树的测试辅助方法，对应 structures/go/TreeNode.go 里的同名函数：
 *
 * <pre>
 * Go                  Java
 * Ints2TreeNode   ->  TreeNodes.build
 * Tree2ints       ->  TreeNodes.toLevelOrder
 * Tree2Preorder   ->  TreeNodes.preorder
 * Tree2Inorder    ->  TreeNodes.inorder
 * Tree2Postorder  ->  TreeNodes.postorder
 * GetTargetNode   ->  TreeNodes.find
 * (tn *TreeNode) Equal -> TreeNodes.equals
 * </pre>
 *
 * <p>Go 版用哨兵 {@code NULL = -1<<63} 表示空节点（{@code []int} 装不了 nil），
 * Java 的 {@code Integer} 能装真正的 null，所以调用写出来和题目给的输入完全一致：
 *
 * <pre>{@code
 * // 题目里的 root = [1,null,2,3]
 * TreeNode root = TreeNodes.build(1, null, 2, 3);
 * }</pre>
 */
public final class TreeNodes {

    private TreeNodes() {}

    /** 按 LeetCode 的层序数组建树，null 表示空节点。 */
    public static TreeNode build(Integer... vals) {
        if (vals.length == 0 || vals[0] == null) {
            return null;
        }

        TreeNode root = new TreeNode(vals[0]);
        Deque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);

        int i = 1;
        while (!queue.isEmpty() && i < vals.length) {
            TreeNode node = queue.poll();
            // 数组只为非空节点列出孩子，所以每出队一个节点就依次消费两个位置
            if (i < vals.length) {
                Integer v = vals[i++];
                if (v != null) {
                    node.left = new TreeNode(v);
                    queue.add(node.left);
                }
            }
            if (i < vals.length) {
                Integer v = vals[i++];
                if (v != null) {
                    node.right = new TreeNode(v);
                    queue.add(node.right);
                }
            }
        }
        return root;
    }

    /**
     * 层序展开，空节点记为 null，末尾的 null 裁掉，便于和题目给的数组直接比对。
     *
     * <p>注意队列里只放非空节点：ArrayDeque 不接受 null 元素，把子节点无脑入队
     * 会直接 NullPointerException。空节点靠"父节点的孩子为空"来记录。
     */
    public static List<Integer> toLevelOrder(TreeNode root) {
        List<Integer> out = new ArrayList<>();
        if (root == null) {
            return out;
        }

        Deque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);
        out.add(root.val);

        while (!queue.isEmpty()) {
            TreeNode node = queue.poll();
            for (TreeNode child : new TreeNode[] {node.left, node.right}) {
                if (child == null) {
                    out.add(null);
                } else {
                    out.add(child.val);
                    queue.add(child);
                }
            }
        }

        while (!out.isEmpty() && out.get(out.size() - 1) == null) {
            out.remove(out.size() - 1);
        }
        return out;
    }

    /** 前序遍历。 */
    public static List<Integer> preorder(TreeNode root) {
        List<Integer> out = new ArrayList<>();
        preorder(root, out);
        return out;
    }

    private static void preorder(TreeNode node, List<Integer> out) {
        if (node == null) {
            return;
        }
        out.add(node.val);
        preorder(node.left, out);
        preorder(node.right, out);
    }

    /** 中序遍历。 */
    public static List<Integer> inorder(TreeNode root) {
        List<Integer> out = new ArrayList<>();
        inorder(root, out);
        return out;
    }

    private static void inorder(TreeNode node, List<Integer> out) {
        if (node == null) {
            return;
        }
        inorder(node.left, out);
        out.add(node.val);
        inorder(node.right, out);
    }

    /** 后序遍历。 */
    public static List<Integer> postorder(TreeNode root) {
        List<Integer> out = new ArrayList<>();
        postorder(root, out);
        return out;
    }

    private static void postorder(TreeNode node, List<Integer> out) {
        if (node == null) {
            return;
        }
        postorder(node.left, out);
        postorder(node.right, out);
        out.add(node.val);
    }

    /** 找出第一个 val 等于 target 的节点，找不到返回 null。层序查找，和 Go 版一致。 */
    public static TreeNode find(TreeNode root, int target) {
        if (root == null) {
            return null;
        }
        Deque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);
        while (!queue.isEmpty()) {
            TreeNode node = queue.poll();
            if (node.val == target) {
                return node;
            }
            // 只放非空节点：ArrayDeque 不接受 null
            if (node.left != null) {
                queue.add(node.left);
            }
            if (node.right != null) {
                queue.add(node.right);
            }
        }
        return null;
    }

    /** 判断两棵树结构和取值是否完全相同。 */
    public static boolean equals(TreeNode a, TreeNode b) {
        if (a == null || b == null) {
            return a == b;
        }
        return a.val == b.val && equals(a.left, b.left) && equals(a.right, b.right);
    }
}
