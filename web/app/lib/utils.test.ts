import { describe, it, expect } from 'vitest';
import { cn, formatTimestamp } from './utils';

describe('cn utility', () => {
  it('should merge class names correctly', () => {
    expect(cn('px-2 py-1', 'text-blue-500')).toBe('px-2 py-1 text-blue-500');
  });

  it('should handle conditional classes', () => {
    expect(cn('base', false && 'conditional', 'always')).toBe('base always');
  });

  it('should merge tailwind classes correctly', () => {
    expect(cn('px-2', 'px-4')).toBe('px-4');
    expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500');
  });
});

describe('formatTimestamp', () => {
  it('should format valid UTC timestamp correctly', () => {
    const timestamp = '2024-01-15T10:30:00.000Z';
    const result = formatTimestamp(timestamp);
    
    expect(result.date).toMatch(/Jan 15, 2024/);
    expect(result.relativeTime).toBeDefined();
    expect(result.fullDate).toBeDefined();
  });

  it('should handle timestamp without Z suffix', () => {
    const timestamp = '2024-01-15T10:30:00.000';
    const result = formatTimestamp(timestamp);
    
    expect(result.date).toMatch(/Jan 15, 2024/);
  });

  it('should return invalid date for malformed timestamp', () => {
    const result = formatTimestamp('invalid-date');
    
    expect(result.date).toBe('Invalid date');
    expect(result.time).toBe('');
    expect(result.fullDate).toBe('');
    expect(result.relativeTime).toBe('');
  });

  it('should calculate relative time correctly', () => {
    const now = new Date();
    const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000);
    const result = formatTimestamp(oneHourAgo.toISOString());
    
    expect(result.relativeTime).toBe('1 hour ago');
  });
});