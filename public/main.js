var app = new Vue({
    el: '#app',
    data: {
        items: [],
        uid: 'default'
    },
    computed: {
        lastItem: function () {
            return this.items ? this.items[this.items.length - 1] : null
        }
    },
    watch: {
        uid: function (newValue) {
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
        onUidSwitchClick: function () {
            var uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
        }
    },
    mounted: function () {
        this.getItems()
    }
})