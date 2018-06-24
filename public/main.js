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
      <div>Current List ID: {{ uid }} </div>
      <input type="text" id="uid-input" placeholder="default" :value="uidInputDisplayValue" @keyup.enter="uidSync">
      <button @click="uidSync">Go</button>
      <div>
        <button @click="generateRandomId">Generate Random ID</button>
      </div>
      <div class="tooltip">
        <button @click="copyToClipboard" @mouseout="showCopiedToClipboard">
          <span class="tooltiptext" id="myTooltip">Copy to Clipboard</span>
          Share
        </button>
      </div>
      <h3>Reminders</h3>
      <div id="listWithHandle">
        <div v-for="item in items">
          <div class="list-group-item">
            <div class="pretty p-default p-thick p-round">
              <input type="checkbox" v-model="item.checked" @change="updateDebounced(item)">
              <div class="state p-success">
                <label></label>
              </div>
            </div>
            <span class="glyphicon glyphicon-move" aria-hidden="true"></span>
            <input type="text" v-model="item.text" class="item-input" :class="{strike: item.checked}" @input="updateDebounced(item)"
              @blur="deleteItemIfEmpty(item)" :ref="item.id">
            <button class="round-btn" :class="{green: item.checked}" @click="deleteItem(item)">X</button>
          </div>
        </div>
      </div>
      <input type="text" v-if="!lastItem || lastItem.text !== ''" placeholder="+ add a reminder" @focus="addNewItem" onfocus="this.placeholder=''"
        onblur="this.placeholder='+ add a reminder'">
    </div>
    `,
    data: function () {
        return {
            items: [],
            uid: getCookie('x-simply-do-uid', 'default')
        }
    },
    computed: {
        lastItem() {
            return this.items ? this.items[this.items.length - 1] : null
        },
        uidInputDisplayValue() {
            return this.uid !== 'default' ? this.uid : ''
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
        },
        uidSync() {
            let uid = document.getElementById('uid-input').value
            this.uid = uid ? uid : 'default'
            api.uid = this.uid
            this.$router.push({
                path: `/list/${this.uid}`
            })
        },
        generateRandomId() {
            this.uid = guid()
        },
        copyToClipboard() {
            const copyText = `https://simply-do.herokuapp.com/${this.$route.path}`
            navigator.clipboard.writeText(copyText).then(function () {
                console.log(`copied ${copyText} to clipboard!`)
            }, function (err) {
                console.error('error copying to clipboard:', err)
            })

            let tooltip = document.getElementById("myTooltip")
            tooltip.innerHTML = "Copied Link to Clipboard"
        },
        showCopiedToClipboard() {
            let tooltip = document.getElementById("myTooltip")
            tooltip.innerHTML = "Copy to Clipboard"
        }
    },
    mounted() {
        let el = document.getElementById('listWithHandle')
        let sortable = Sortable.create(el)
        Sortable.create(listWithHandle, {
            handle: '.glyphicon-move',
            animation: 150
        })

        if (this.$route.params.id) {
            this.uid = this.$route.params.id
            setCookie('x-simply-do-uid', this.uid, 1)
        }
        api.uid = this.uid
        api.getItems()
            .then(items => this.items = items)
    },
    beforeRouteUpdate(to, from, next) {
        console.debug('to', to)
        console.debug('from', from)
        this.uid = to.params.id
        setCookie('x-simply-do-uid', this.uid, 1)
        api.uid = this.uid
        api.getItems()
            .then(items => this.items = items)
        next()
        Sortable.create(listWithHandle, {
            handle: '.glyphicon-move',
            animation: 150
        });
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
    mounted() {
        if (this.$route.path === '/') {
            this.$router.push(`list/${getCookie('x-simply-do-uid', 'default')}`)
        }
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

function guid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1)
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4()
}