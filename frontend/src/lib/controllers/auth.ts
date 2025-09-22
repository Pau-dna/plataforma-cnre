import { BaseController } from './base';
import type {
    Content,
    CreateContentDTO,
    UpdateContentDTO,
    ReorderItemDTO,
    ModuleContent,
    User
} from '$lib/types';
import type { AuthTokens } from '$lib/types/tokens';
import type { SignInResponse } from '$lib/types/dto/auth';

export class AuthController extends BaseController {
    
    async loginWithGoogle(code: string) {
        return this.post<SignInResponse>('/auth/google', { code });
    }

    async getUserInfo() {
        return this.get<User>('/auth/me');
    }

}
