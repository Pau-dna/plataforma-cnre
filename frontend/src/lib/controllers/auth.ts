import { BaseController } from './base';
import type {
    Content,
    CreateContentDTO,
    UpdateContentDTO,
    ReorderItemDTO,
    ModuleContent
} from '$lib/types';
import type { AuthTokens } from '$lib/types/tokens';

export class AuthController extends BaseController {
    
    async loginWithGoogle(code: string) {
        return this.post<AuthTokens>('/auth/google', { code });
    }

}
