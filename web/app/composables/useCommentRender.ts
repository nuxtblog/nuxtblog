/**
 * Renders comment content with @mention highlighting.
 * HTML-escapes the raw text first to prevent XSS, then wraps @mentions in a styled span.
 */
export const useCommentRender = () => {
  const escapeHtml = (text: string) =>
    text
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");

  /**
   * Returns safe HTML string with @mentions highlighted.
   * Usage: <p v-html="renderComment(content)" />
   */
  const renderComment = (text: string): string => {
    if (!text) return "";
    const escaped = escapeHtml(text);
    // Match @mention at start or after whitespace, ends at whitespace or end-of-string
    return escaped.replace(
      /(^|\s)(@[^\s@<]{1,50})/g,
      '$1<span class="text-info font-medium">$2</span>',
    );
  };

  return { renderComment };
};
