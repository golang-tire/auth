import axios from 'axios';
import {configs} from 'services/Network/config';
import TokenService from "./tokenService";
import history from "../../history";

const BASE_AUTH_URL = configs.API_URL + "/auth"

axios.interceptors.request.use(req => {
    req.headers.authorization = TokenService.getAuthentication().headers.Authorization;
    req.headers.contentType = "application/json"
    return req;
});

axios.interceptors.response.use( (response) => {
    // Return a successful response back to the calling service
    return response;
}, (error) => {
    // Return any error which is not due to authentication back to the calling service
    if (error.response.status !== 401) {
        return new Promise((resolve, reject) => {
            reject(error);
        });
    }

    if (error.response.status === 401 && error.config.url === BASE_AUTH_URL + '/login') {
        return new Promise((resolve, reject) => {
            reject(error);
        });
    }

    // Logout user if token refresh didn't work or user is disabled
    if (error.config.url === BASE_AUTH_URL + '/token/refresh' || error.response.status === 500) {
        TokenService.clear();
        history.push("/login")
        return new Promise((resolve, reject) => {
            reject(error);
        });
    }
    // Try request again with new token
    return TokenService.getNewToken()
        .then((token) => {
            // New request with new token
            const config = error.config;
            config.headers['authorization'] = `${configs.BEARER_PARAM} ${token}`;
            return new Promise((resolve, reject) => {
                axios.request(config).then(response => {
                    resolve(response);
                }).catch((error) => {
                    reject(error);
                })
            });
        })
        .catch((error) => {
            Promise.reject(error);
        });
});

// AuthService service is responsible for auth login, logout
const AuthService = {
    login(username, password) {
        return axios.post(
            BASE_AUTH_URL + "/login",
            {
                username,
                password
            }
        ).then((res) =>{
            if (res.status === 200) {
                TokenService.storeToken(res.data.access_token);
                TokenService.storeRefreshToken(res.data.refresh_token);
            }
            return {
                status: res.status,
                ...res.data
            }
        },(error) => {
            return {
                status: error.response.status,
                error: error
            }
        })
    },
    forgotPassword(email){
        return axios.post(
            BASE_AUTH_URL + "/forgot-password",
            {email: email}
        ).then((res) =>{
            return {
                status: res.status,
                ...res.data
            }
        },(error) => {
            return {
                status: error.response.status,
                error: error
            }
        })
    },
    logout() {
        return axios.post(
            BASE_AUTH_URL + "/logout",
        ).then((res) =>{
            TokenService.clear();
            return {
                status: res.status,
                ...res.data
            }
        },(error) => {
            return {
                status: error.response.status,
                error: error
            }
        })
    },
};

export default AuthService;