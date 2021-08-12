import axios from 'axios'


const baseURL = import.meta.env.MODE === 'development' ? '' : window.location.protocol + "//" + window.location.host

const _request = axios.create({
    baseURL: baseURL,
    timeout: 10000,
})

_request.interceptors.request.use(
    (config) => {
        console.log(config.url)
        return config
    },
    (error) => {
        return error
    }
)

_request.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        return error
    }
)

export const get = (url, params = {}) => {
    return new Promise((resolve, reject) => {
        _request({
            url: url,
            method: 'GET',
            params: params,
        }).then((res) => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}

export const post = (url, params={}, headers={
    'Content-Type': 'application/json'
}) => {
    return new Promise((resolve, reject) => {
        _request({
            url: url,
            method: 'POST',
            data: params,
            headers: headers
        }).then(res=>{
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}