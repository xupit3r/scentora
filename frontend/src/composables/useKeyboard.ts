/**
 * Keyboard Shortcuts Composable
 * Provides keyboard shortcut functionality
 */

import { onMounted, onUnmounted } from 'vue';

export interface KeyboardShortcut {
  key: string;
  ctrl?: boolean;
  meta?: boolean;
  shift?: boolean;
  alt?: boolean;
  handler: () => void;
  description: string;
}

export function useKeyboard(shortcuts: KeyboardShortcut[]) {
  function handleKeydown(event: KeyboardEvent) {
    for (const shortcut of shortcuts) {
      const ctrlMatch = shortcut.ctrl ? event.ctrlKey : !event.ctrlKey;
      const metaMatch = shortcut.meta ? event.metaKey : !event.metaKey;
      const shiftMatch = shortcut.shift ? event.shiftKey : !event.shiftKey;
      const altMatch = shortcut.alt ? event.altKey : !event.altKey;
      
      if (
        event.key.toLowerCase() === shortcut.key.toLowerCase() &&
        (ctrlMatch || metaMatch) &&
        shiftMatch &&
        altMatch
      ) {
        event.preventDefault();
        shortcut.handler();
        break;
      }
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown);
  });

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown);
  });

  return {
    shortcuts,
  };
}

export function useSimpleKeyboard(key: string, handler: () => void, enabled = true) {
  function handleKeydown(event: KeyboardEvent) {
    if (!enabled) return;
    
    // Don't trigger if user is typing in an input
    const target = event.target as HTMLElement;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') {
      return;
    }

    if (event.key === key) {
      event.preventDefault();
      handler();
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown);
  });

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown);
  });
}
