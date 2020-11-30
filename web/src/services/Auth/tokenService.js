import axios from 'axios';
import {configs} from 'services/Network/config';

const LOCAL_STORAGE_TOKEN = 'access_token';
const LOCAL_STORAGE_REFRESH_TOKEN = 'refresh_token';
const REFRESH_TOKEN_URL = configs.API_URL + "/auth/token/refresh"

// TokenService service is responsible for session storage
const TokenService = {
    isAuthenticated(){
        return TokenService.getToken() !== null
    },
    getAuthentication() {
        return {
            headers: { 'Authorization': configs.BEARER_PARAM + ' ' + this.getToken() }
        };
    },
    getNewToken() {
        return new Promise((resolve, reject) => {
            axios
                .post(REFRESH_TOKEN_URL, { refresh_token: TokenService.getRefreshToken() })
                .then(response => {
                    TokenService.storeToken(response.data.access_token);
                    TokenService.storeRefreshToken(response.data.refresh_token);
                    resolve(response.data.access_token);
                })
                .catch((error) => {
                    reject(error);
                });
        });
    },
    storeToken(token){
        localStorage.setItem(LOCAL_STORAGE_TOKEN, token);
    },
    storeRefreshToken(refreshToken) {
        localStorage.setItem(LOCAL_STORAGE_REFRESH_TOKEN, refreshToken);
    },
    clear(){
        localStorage.removeItem(LOCAL_STORAGE_TOKEN);
        localStorage.removeItem(LOCAL_STORAGE_REFRESH_TOKEN);
    },
    getRefreshToken() {
        return localStorage.getItem(LOCAL_STORAGE_REFRESH_TOKEN);
    },
    getToken(){
        return localStorage.getItem(LOCAL_STORAGE_TOKEN);
    }
};

export default TokenService;