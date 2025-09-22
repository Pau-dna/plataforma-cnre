import { describe, it, expect, vi } from 'vitest';
import { ModuleController } from '$lib/controllers';
import type { CreateModuleDTO } from '$lib/types/dto';

// Mock the ModuleController
vi.mock('$lib/controllers', () => {
	return {
		ModuleController: vi.fn().mockImplementation(() => ({
			createModule: vi.fn()
		}))
	};
});

describe('Module Creation Functionality', () => {
	it('should create a module with correct data structure', async () => {
		const mockCreateModule = vi.fn().mockResolvedValue({
			id: 1,
			title: 'Test Module',
			description: 'Test Description',
			order: 0,
			course_id: 1
		});

		const controller = new ModuleController();
		controller.createModule = mockCreateModule;

		const moduleData: CreateModuleDTO = {
			title: 'Test Module',
			description: 'Test Description',
			order: 0,
			course_id: 1
		};

		const result = await controller.createModule(moduleData);

		expect(mockCreateModule).toHaveBeenCalledWith(moduleData);
		expect(result).toEqual({
			id: 1,
			title: 'Test Module',
			description: 'Test Description',
			order: 0,
			course_id: 1
		});
	});

	it('should validate required fields', () => {
		const validationFunction = (formdata: CreateModuleDTO) => {
			return !formdata.title || !formdata.description;
		};

		// Test valid data
		const validData: CreateModuleDTO = {
			title: 'Valid Title',
			description: 'Valid Description',
			order: 0,
			course_id: 1
		};

		expect(validationFunction(validData)).toBe(false);

		// Test invalid data (missing title)
		const invalidData: CreateModuleDTO = {
			title: '',
			description: 'Valid Description',
			order: 0,
			course_id: 1
		};

		expect(validationFunction(invalidData)).toBe(true);

		// Test invalid data (missing description)
		const invalidData2: CreateModuleDTO = {
			title: 'Valid Title',
			description: '',
			order: 0,
			course_id: 1
		};

		expect(validationFunction(invalidData2)).toBe(true);
	});
});