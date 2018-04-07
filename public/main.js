var app = new Vue({
    el: '#app',
    data: {
        items: []
    },
    computed: {
        lastItem: function () {
            return this.items ? this.items[this.items.length - 1] : null
        }
    },
    methods: {
        getItems: function () {
            return fetch('/api/items').then(function (response) {
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
                body: JSON.stringify(item)
            }).then(function (response) {
                return response.json()
            })

        },
        deleteItem: function (item) {
            fetch('/api/items/' + item.id, {
                method: 'DELETE'
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
        }
    },
    mounted: function () {
        this.getItems()
    }
})