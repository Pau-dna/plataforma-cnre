import { describe, it, expect } from 'vitest';
import { detectUrls, normalizeUrl } from './url';

describe('URL Detection', () => {
	it('should detect HTTP URLs', () => {
		const text = 'Visit https://example.com for more info';
		const result = detectUrls(text);

		expect(result).toHaveLength(3);
		expect(result[0]).toEqual({ text: 'Visit ', isUrl: false });
		expect(result[1]).toEqual({ text: 'https://example.com', isUrl: true });
		expect(result[2]).toEqual({ text: ' for more info', isUrl: false });
	});

	it('should detect www URLs', () => {
		const text = 'Check www.example.com today';
		const result = detectUrls(text);

		expect(result).toHaveLength(3);
		expect(result[0]).toEqual({ text: 'Check ', isUrl: false });
		expect(result[1]).toEqual({ text: 'www.example.com', isUrl: true });
		expect(result[2]).toEqual({ text: ' today', isUrl: false });
	});

	it('should handle text without URLs', () => {
		const text = 'No links here';
		const result = detectUrls(text);

		expect(result).toHaveLength(1);
		expect(result[0]).toEqual({ text: 'No links here', isUrl: false });
	});

	it('should handle multiple URLs', () => {
		const text = 'Visit https://google.com and www.github.com';
		const result = detectUrls(text);

		expect(result).toHaveLength(4);
		expect(result[0]).toEqual({ text: 'Visit ', isUrl: false });
		expect(result[1]).toEqual({ text: 'https://google.com', isUrl: true });
		expect(result[2]).toEqual({ text: ' and ', isUrl: false });
		expect(result[3]).toEqual({ text: 'www.github.com', isUrl: true });
	});

	it('should normalize www URLs', () => {
		expect(normalizeUrl('www.example.com')).toBe('https://www.example.com');
		expect(normalizeUrl('https://example.com')).toBe('https://example.com');
	});
});
