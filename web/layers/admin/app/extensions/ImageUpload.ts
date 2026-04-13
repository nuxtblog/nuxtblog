import { Node, type CommandProps } from "@tiptap/core";
import { VueNodeViewRenderer } from "@tiptap/vue-3";
import EditorImageUploadNode from "../components/EditorImageUploadNode.vue";

declare module "@tiptap/core" {
  interface Commands<ReturnType> {
    imageUpload: {
      setImageUpload: () => ReturnType;
    };
  }
}

export const ImageUpload = Node.create({
  name: "imageUpload",
  group: "block",
  atom: true,
  draggable: true,

  addAttributes() {
    return {};
  },

  parseHTML() {
    return [];
  },

  renderHTML() {
    return ["div", { "data-type": "image-upload" }];
  },

  addNodeView() {
    return VueNodeViewRenderer(EditorImageUploadNode);
  },

  addCommands() {
    return {
      setImageUpload:
        () =>
        ({ commands }: CommandProps) => {
          return commands.insertContent({ type: this.name });
        },
    };
  },
});
