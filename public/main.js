class API {
    constructor() {
        this.uid = getCookie('x-simply-do-uid', 'default')
    }

    getItems() {
        console.debug('GET /api/items', this.uid)
        return fetch('/api/items', {
                headers: this._default_headers
            })
            .then(res => res.json())
            .catch(e => console.error(e))
    }

    createItem() {
        console.debug('POST /api/items', this.uid)
        return fetch('/api/items', {
                method: 'POST',
                headers: this._default_headers,
                body: JSON.stringify({
                    checked: false,
                    text: ''
                })
            })
            .then(res => res.json())
            .catch(e => console.error(e))
    }

    updateItem(item) {
        console.debug('PUT /api/items', this.uid)
        return fetch('/api/items/' + item.id, {
                method: 'PUT',
                headers: this._default_headers,
                body: JSON.stringify(item)
            })
            .then(res => res.json())
            .catch(e => console.error(e))
    }

    deleteItem(item) {
        console.debug('DELETE /api/items', this.uid, item)
        return fetch('/api/items/' + item.id, {
                method: 'DELETE',
                headers: this._default_headers
            })
            .then(res => res.json())
            .catch(e => console.error(e))
    }

    get _default_headers() {
        return {
            'x-simply-do-uid': this.uid
        }
    }
}

class WorkQueue {
    constructor(api) {
        this.api = api
        this.workMap = {}
    }

    enqueueUpdate(item) {
        clearTimeout(this.workMap[item.id])
        this.workMap[item.id] = setTimeout(() => {
            if (item.text !== '') {
                api.updateItem(item).then(item => {
                    console.debug("UPDATE finished:", item)
                    delete this.workMap[item.id]
                })
            }
        }, 1000)
    }
}

api = new API()
workQueue = new WorkQueue(api)

const List = {
    template: `
    <div>
        <div v-for="item in items">
        <div class="pretty p-default p-thick p-round">
            <input type="checkbox" v-model="item.checked" @change="updateDebounced(item)">
            <div class="state p-success">
            <label></label>
            </div>
        </div>
        <input type="text" v-model="item.text" class="item-input" :class="{strike: item.checked}" @input="updateDebounced(item)"
            @blur="deleteItemIfEmpty(item)" :ref="item.id">
        <button class="round-btn" :class="{green: item.checked}" @click="deleteItem(item)">X</button>
        </div>
        <input type="text" v-if="!lastItem || lastItem.text !== ''" placeholder="+ add a reminder" @focus="addNewItem" onfocus="this.placeholder=''"
        onblur="this.placeholder='+ add a reminder'">
    </div>
    `,
    data: function () {
        return {
            items: []
        }
    },
    computed: {
        lastItem() {
            return this.items ? this.items[this.items.length - 1] : null
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

            Vue.nextTick()
                .then(() => this.$refs[this.lastItem.id][0].focus())

            api.createItem()
                .then(item => {
                    newItem.id = item.id
                    newItem.time_created = item.time_created
                    console.debug("CREATE finished:", newItem.id, newItem)
                })
                .catch(e => console.error("CREATE errored:", e))
        },
        updateDebounced(item) {
            workQueue.enqueueUpdate(item)
        },
        deleteItem(item) {
            this._executeAsyncApiDeleteItem(item)
            this.items.splice(this.items.indexOf(item), 1)
        },
        deleteItemIfEmpty(item) {
            if (item.text === '') {
                this.deleteItem(item)
            }
        },
        _executeAsyncApiDeleteItem(item) {
            setTimeout(() => {
                if (item.id === 'id-placeholder') {
                    console.debug('DELETE id-placeholder: delaying for another 1 second...')
                    this._executeAsyncApiDeleteItem(item)
                    return
                }
                api.deleteItem(item)
                    .then(() => console.debug("DELETE finished:", item.id, item))
            }, 1000)
        }
    },
    mounted() {
        api.getItems()
            .then(items => this.items = items)
    },
    beforeRouteUpdate(to, from, next) {
        console.debug('to', to)
        console.debug('from', from)
        api.uid = to.params.id
        api.getItems()
            .then(items => this.items = items)
        next()
    }
}

const router = new VueRouter({
    routes: [{
        path: '/list/:id',
        component: List
    }]
})

const app = new Vue({
    router,
    data: {
        uid: getCookie('x-simply-do-uid', 'default')
    },
    computed: {
        uidInputDisplayValue() {
            return this.uid !== 'default' ? this.uid : ''
        }
    },
    watch: {
        uid(newValue) {
            setCookie('x-simply-do-uid', newValue, 1)
            api.uid = newValue
        }
    },
    methods: {
        uidSync() {
            let uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
            this.$router.push({
                path: `/list/${this.uid}`
            })
        },
    },
    mounted() {
        this.$router.push({
            path: `/list/${this.uid}`
        })
    }
}).$mount('#app')

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
        let c = ca[i]
        while (c.charAt(0) == ' ') {
            c = c.substring(1)
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length)
        }
    }
    return defaultvalue
}