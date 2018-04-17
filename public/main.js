var api = {
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

var workQueue = {
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

var app = new Vue({
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
            var newItem = {
                checked: false,
                text: '',
                uniquePlaceholderId: guid()
            }
            this.items.push(newItem)
            Vue.nextTick().then(function () {
                app.$refs[app.lastItem.id][0].focus()
            })
            api.createItem().then(function (item) {
                if (app._shouldClean(newItem.uniquePlaceholderId)) {
                    console.debug("CREATE finished: should clean...")
                    api.deleteItem(item).then(function () {
                        console.debug("DELETE finished:", item.id, item)
                    })
                    return
                }
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
            if (item.id) {
                api.deleteItem(item).then(function () {
                    console.debug("DELETE finished:", item.id, item)
                })
            }
            this.items.splice(this.items.indexOf(item), 1)
        },
        deleteItemIfEmpty: function (item) {
            if (item.text === '') {
                this.deleteItem(item)
            }
        },
        uidSync: function () {
            var uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
        },
        _shouldClean(uniquePlaceholderId) {
            for (var i = 0; i < this.items.length; i++) {
                if (this.items[i].uniquePlaceholderId === uniquePlaceholderId) {
                    return false
                }
            }
            return true
        }
    },
    mounted: function () {
        api.getItems().then(function (items) {
            app.items = items
        })
    }
})

function setCookie(cname, cvalue, exdays) {
    var d = new Date()
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000))
    var expires = 'expires=' + d.toUTCString()
    document.cookie = cname + '=' + cvalue + ';' + expires + ';path=/'
}

function getCookie(cname, defaultvalue) {
    var name = cname + '='
    var decodedCookie = decodeURIComponent(document.cookie)
    var ca = decodedCookie.split(';')
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1)
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length)
        }
    }
    return defaultvalue
}

function guid() {
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
        s4() + '-' + s4() + s4() + s4();
}

function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
}