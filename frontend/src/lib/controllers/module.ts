import type { Module } from "$lib/types/models/course";

export class ModuleController {
    async getModules(courseID :number): Promise<Module[]> {
        return [
            {
                id: 1,
                title: 'Module 1',
                description: 'Description for Module 1',
            }
        ];
    }
}
