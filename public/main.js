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
            this.getItems()
        }
    },
    methods: {
        getItems: function () {
            return fetch('/api/items', {
                headers: {
                    'x-simply-do-uid': this.uid
                }
            }).then(function (response) {
                return response.json()
            }).then(function (body) {
                console.log(body)
                app.items = body
                return body
            })
        },
        createItem: function () {
            var item = {
                checked: false,
                text: ''
            }
            return fetch('/api/items', {
                method: 'POST',
                headers: {
                    'x-simply-do-uid': this.uid
                },
                body: JSON.stringify(item)
            }).then(function (response) {
                return response.json()
            }).then(function (body) {
                console.log(body)
            }).then(function () {
                app.getItems().then(function () {
                    app.$refs[app.lastItem.id][0].focus()
                })
            })
        },
        updateItem: function (item) {
            return fetch('/api/items/' + item.id, {
                method: 'PUT',
                headers: {
                    'x-simply-do-uid': this.uid
                },
                body: JSON.stringify(item)
            }).then(function (response) {
                return response.json()
            })

        },
        deleteItem: function (item) {
            fetch('/api/items/' + item.id, {
                method: 'DELETE',
                headers: {
                    'x-simply-do-uid': this.uid
                }
            }).then(function (response) {
                return response.json()
            }).then(function (body) {
                app.items = body
                return body
            })
        },
        deleteItemIfEmpty: function (item) {
            if (item.text === '') {
                this.deleteItem(item)
            }
        },
        uidSync: function () {
            var uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
        }
    },
    mounted: function () {
        this.getItems()
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