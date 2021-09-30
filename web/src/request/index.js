import axios from 'axios'
import cogoToast from 'cogo-toast'

const baseURL = import.meta.env.MODE === 'development' ? '' : window.location.protocol + "//" + window.location.host

const _request = axios.create({
    baseURL: baseURL,
    timeout: 10000,
})

const cogoError = (msg) => {
    cogoToast.error(msg, {
        heading: 'Request failed'
    })
}

// send request
_request.interceptors.request.use(
    (config) => {
        console.log(config.url)
        return config
    },
    (error) => {
        return error
    }
)

// response
_request.interceptors.response.use(
    (response) => {
        // check response's code
        if (response.status === 200) {
            const res = response.data
            switch (res.code) {
                case 0:
                    break
                default:
                    cogoError('Call api failed...')
                    break
            }
            return res
        }
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

export const post = (url, params = {}, headers = {
    'Content-Type': 'application/json'
}) => {
    return new Promise((resolve, reject) => {
        _request({
            url: url,
            method: 'POST',
            data: params,
            headers: headers
        }).then(res => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}

export const del = (url, params = {}, headers = {
    'Content-Type': 'application/json'
}) => {
    return new Promise((resolve, reject) => {
        _request({
            url: url,
            method: 'DELETE',
            params: params,
            headers: headers
        }).then(res => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}

export const update = (url, params = {}, headers = {
    'Content-Type': 'application/json'
}) => {
    return new Promise((resolve, reject) => {
        _request({
            url: url,
            method: 'PATCH',
            data: params,
            headers: headers
        }).then(res => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}