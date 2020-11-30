import axios from 'axios';
import {configs} from 'services/Network/config';

// ApiService service is responsible for call rest api
const ApiService = {
    getAll(Url, args) {
        let url = configs.API_URL + "/" + Url
        return axios.get(url, {params: {...args}})
    },
    get(Url, uuid) {
        const url = configs.API_URL + "/" + Url + "/" + uuid;
        return axios.get(url)
    },
    post(Url, data) {
        const url = configs.API_URL + "/" + Url;
        return axios.post(url, data)
    },
    put(Url, Uuid, data) {
        const url = configs.API_URL + "/" + Url + "/" + Uuid;
        return axios.put(url, data)
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