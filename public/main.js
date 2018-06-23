let api = {
    uid: getCookie('x-simply-do-uid', 'default'),
    getItems: function () {
        console.debug('GET /api/items', this.uid)
        return fetch('/api/items', {
            headers: this._default_headers()
        }).then(function (response) {
            return response.json()
        }).catch(function (e) {
            console.error(e)
        })
    },
    createItem: function () {
        console.debug('POST /api/items', this.uid)
        return fetch('/api/items', {
            method: 'POST',
            headers: this._default_headers(),
            body: JSON.stringify({
                checked: false,
                text: ''
            })
        }).then(function (response) {
            return response.json()
        }).catch(function (e) {
            console.error(e)
        })
    },
    updateItem: function (item) {
        console.debug('PUT /api/items', this.uid)
        return fetch('/api/items/' + item.id, {
            method: 'PUT',
            headers: this._default_headers(),
            body: JSON.stringify(item)
        }).then(function (response) {
            return response.json()
        }).catch(function (e) {
            console.error(e)
        })
    },
    deleteItem: function (item) {
        console.debug('DELETE /api/items', this.uid, item)
        return fetch('/api/items/' + item.id, {
            method: 'DELETE',
            headers: this._default_headers()
        }).then(function (response) {
            return response.json()
        }).catch(function (e) {
            console.error(e)
        })
    },
    _default_headers: function () {
        return {
            'x-simply-do-uid': this.uid
        }
    },
}

let workQueue = {
    addUpdate: function (item) {
        clearTimeout(this[item.id])
        this[item.id] = setTimeout(function () {
            if (item.text !== '') {
                api.updateItem(item).then(function (item) {
                    console.debug("UPDATE finished:", item)
                    delete workQueue[item.id]
                })
            }
        }, 1000)
    }
}

let app = new Vue({
    el: '#app',
    data: {
        items: [],
        uid: getCookie('x-simply-do-uid', 'default')
    },
    computed: {
        uidInputDisplayValue: function () {
            return this.uid !== 'default' ? this.uid : ''
        },
        lastItem: function () {
            return this.items ? this.items[this.items.length - 1] : null
        }
    },
    watch: {
        uid: function (newValue) {
            setCookie('x-simply-do-uid', newValue, 1)
            api.uid = newValue
            api.getItems().then(function (items) {
                app.items = items
            })
        }
    },
    methods: {
        addNewItem() {
            let newItem = {
                id: 'id-placeholder',
                checked: false,
                text: '',
                time_created: 'time_created-placeholder'
            }
            this.items.push(newItem)
            Vue.nextTick().then(function () {
                app.$refs[app.lastItem.id][0].focus()
            })
            api.createItem().then(function (item) {
                newItem.id = item.id
                newItem.time_created = item.time_created
                console.debug("CREATE finished: updated new item:", newItem.id, newItem)
            }).catch(function (e) {
                console.error("CREATE errored:", e)
            })
        },
        updateDebounced: function (item) {
            workQueue.addUpdate(item)
        },
        deleteItem: function (item) {
            this._executeAsyncApiDeleteItem(item)
            this.items.splice(this.items.indexOf(item), 1)
        },
        deleteItemIfEmpty: function (item) {
            if (item.text === '') {
                this.deleteItem(item)
            }
        },
        uidSync: function () {
            let uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
        },
        _executeAsyncApiDeleteItem(item) {
            setTimeout(function () {
                if (item.id === 'id-placeholder') {
                    console.debug('DELETE id-placeholder: delaying...')
                    app._executeApiDeleteItem(item)
                    return
                }
                api.deleteItem(item).then(function () {
                    console.debug("DELETE finished:", item.id, item)
                })
            }, 1000)
        }
    },
    mounted: function () {
        api.getItems().then(function (items) {
            app.items = items
        })
    }
})

function setCookie(cname, cvalue, exdays) {
    let d = new Date()
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000))
    let expires = 'expires=' + d.toUTCString()
    document.cookie = cname + '=' + cvalue + ';' + expires + ';path=/'
}

function getCookie(cname, defaultvalue) {
    let name = cname + '='
    let decodedCookie = decodeURIComponent(document.cookie)
    let ca = decodedCookie.split(';')
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1)
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length)
        }
    }
    return defaultvalue
}