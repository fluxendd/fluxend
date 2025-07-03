import { describe, it, expect, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useToast } from './use-toast';

describe('useToast Hook', () => {
  it('should start with no toasts', () => {
    const { result } = renderHook(() => useToast());
    expect(result.current.toasts).toHaveLength(0);
  });

  it('should add a toast', () => {
    const { result } = renderHook(() => useToast());
    
    act(() => {
      result.current.toast({
        title: 'Test Toast',
        description: 'This is a test',
      });
    });

    expect(result.current.toasts).toHaveLength(1);
    expect(result.current.toasts[0]).toMatchObject({
      title: 'Test Toast',
      description: 'This is a test',
    });
  });

  it('should dismiss a toast', () => {
    const { result } = renderHook(() => useToast());
    
    let toastId: string;
    act(() => {
      const { id } = result.current.toast({
        title: 'Test Toast',
      });
      toastId = id;
    });

    expect(result.current.toasts).toHaveLength(1);

    act(() => {
      result.current.dismiss(toastId!);
    });

    // Toast should be marked as dismissed but still in the array
    expect(result.current.toasts[0].open).toBe(false);
  });

  it('should handle multiple toasts', () => {
    const { result } = renderHook(() => useToast());
    
    act(() => {
      result.current.toast({ title: 'Toast 1' });
      result.current.toast({ title: 'Toast 2' });
      result.current.toast({ title: 'Toast 3' });
    });

    expect(result.current.toasts).toHaveLength(3);
    expect(result.current.toasts[0].title).toBe('Toast 1');
    expect(result.current.toasts[1].title).toBe('Toast 2');
    expect(result.current.toasts[2].title).toBe('Toast 3');
  });

  it('should respect toast limit', () => {
    const { result } = renderHook(() => useToast());
    
    // Add more toasts than the limit
    act(() => {
      for (let i = 0; i < 10; i++) {
        result.current.toast({ title: `Toast ${i}` });
      }
    });

    // Should only keep the last TOAST_LIMIT toasts
    expect(result.current.toasts.length).toBeLessThanOrEqual(5); // Assuming TOAST_LIMIT is 5
  });
});