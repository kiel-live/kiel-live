import { describe, expect, it } from 'vitest';
import { get24hTime } from './date';

describe('get24hTime', () => {
  it('formats a local datetime', () => {
    expect(get24hTime('2026-10-30T09:05:00')).toBe('09:05');
  });

  it('returns an empty string for a missing time', () => {
    expect(get24hTime('')).toBe('');
  });

  it('returns an empty string for an invalid time', () => {
    expect(get24hTime('not-a-date')).toBe('');
  });
});
