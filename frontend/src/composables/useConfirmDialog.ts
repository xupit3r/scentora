import { useDialog } from 'naive-ui';

/**
 * Composable for showing confirmation dialogs
 * Uses Naive UI's dialog component for consistent UX
 */
export function useConfirmDialog() {
  const dialog = useDialog();

  /**
   * Show a confirmation dialog
   * @param options Dialog options
   * @returns Promise that resolves to true if confirmed, false if cancelled
   */
  function confirm(options: {
    title: string;
    content: string;
    positiveText?: string;
    negativeText?: string;
    type?: 'error' | 'warning' | 'info' | 'success';
  }): Promise<boolean> {
    return new Promise((resolve) => {
      dialog[options.type || 'warning']({
        title: options.title,
        content: options.content,
        positiveText: options.positiveText || 'Confirm',
        negativeText: options.negativeText || 'Cancel',
        onPositiveClick: () => {
          resolve(true);
        },
        onNegativeClick: () => {
          resolve(false);
        },
        onClose: () => {
          resolve(false);
        },
      });
    });
  }

  /**
   * Show a delete confirmation dialog
   * @param itemName Name of the item being deleted
   * @returns Promise that resolves to true if confirmed, false if cancelled
   */
  async function confirmDelete(itemName: string): Promise<boolean> {
    return confirm({
      title: 'Confirm Deletion',
      content: `Are you sure you want to delete "${itemName}"? This action cannot be undone.`,
      positiveText: 'Delete',
      negativeText: 'Cancel',
      type: 'error',
    });
  }

  /**
   * Show a warning dialog (for potentially destructive actions)
   * @param title Dialog title
   * @param message Warning message
   * @returns Promise that resolves to true if confirmed, false if cancelled
   */
  async function confirmAction(title: string, message: string): Promise<boolean> {
    return confirm({
      title,
      content: message,
      positiveText: 'Continue',
      negativeText: 'Cancel',
      type: 'warning',
    });
  }

  return {
    confirm,
    confirmDelete,
    confirmAction,
  };
}
