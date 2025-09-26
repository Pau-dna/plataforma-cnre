/**
 * Utility functions for URL detection and processing
 */

/**
 * Regular expression to detect URLs in text
 * Matches http://, https://, and www. patterns
 */
const URL_REGEX = /(https?:\/\/[^\s]+)|(www\.[^\s]+)/gi;

/**
 * Detects URLs in text and returns an array of text segments and URLs
 * @param text The input text to process
 * @returns Array of objects with {text, isUrl} properties
 */
export function detectUrls(text: string): Array<{ text: string; isUrl: boolean }> {
	if (!text) return [];
	
	const result: Array<{ text: string; isUrl: boolean }> = [];
	let lastIndex = 0;
	
	// Find all URL matches
	const matches = [...text.matchAll(URL_REGEX)];
	
	for (const match of matches) {
		const matchStart = match.index!;
		const matchEnd = matchStart + match[0].length;
		
		// Add text before the URL
		if (matchStart > lastIndex) {
			const beforeText = text.slice(lastIndex, matchStart);
			if (beforeText) {
				result.push({ text: beforeText, isUrl: false });
			}
		}
		
		// Add the URL
		result.push({ text: match[0], isUrl: true });
		lastIndex = matchEnd;
	}
	
	// Add remaining text after the last URL
	if (lastIndex < text.length) {
		const remainingText = text.slice(lastIndex);
		if (remainingText) {
			result.push({ text: remainingText, isUrl: false });
		}
	}
	
	// If no URLs found, return the original text
	if (result.length === 0) {
		result.push({ text, isUrl: false });
	}
	
	return result;
}

/**
 * Ensures a URL has a proper protocol
 * @param url The URL to normalize
 * @returns URL with proper protocol
 */
export function normalizeUrl(url: string): string {
	if (url.startsWith('www.')) {
		return `https://${url}`;
	}
	return url;
}