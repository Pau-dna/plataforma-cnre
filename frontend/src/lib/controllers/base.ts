export class BaseController {
    token?: string;

    constructor(token?: string) {
        this.token = token;
    }
}