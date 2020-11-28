import axios from 'axios';
import {configs} from 'services/Network/config';
import AuthService from "services/Auth/authService";
import history from "../../history";


const isEmpty = (obj) => {
    for(let key in obj) {
        if(obj.hasOwnProperty(key))
            return false;
    }
    return true;
}

axios.interceptors.request.use(req => {
    const headers = AuthService.getAuthHeader()
    if (isEmpty(headers)) {
        history.push("/login")
        return
    }
    req.headers.authorization = headers.Authorization;
    return req;
});

axios.interceptors.response.use(
    res => res,
    err => {
        if (err.response.status === 401) {
            history.push("/login")
            throw new Error(`Session expired`);
        }
        throw err;
    }
);

// ApiService service is responsible for call rest api
const ApiService = {
    getAll(Url) {
        const url = configs.API_URL + "/" + Url
        return axios.get(url)
    },
    get(Url, uuid) {
        const url = configs.API_URL + "/" + Url + "/" + uuid;
        return axios.get(url)
    },
    delete(Url, uuid) {
        let url = configs.API_URL + "/" + Url + "/" + uuid;
        return axios.delete(url)
    },
    request(method, url, data){
        return axios({method, url, data})
    }
};

export default ApiService;