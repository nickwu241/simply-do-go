<template>
  <div>
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
          <input type="text" v-model="item.text" class="item-input" :class="{strike: item.checked}" @input="updateDebounced(item)" @blur="deleteItemIfEmpty(item)" :ref="item.id">
          <button class="round-btn" :class="{green: item.checked}" @click="deleteItem(item)">X</button>
        </div>
      </div>
    </div>
    <input type="text" v-if="!lastItem || lastItem.text !== ''" placeholder="+ add a reminder" @focus="addNewItem" onfocus="this.placeholder=''" onblur="this.placeholder='+ add a reminder'">
  </div>
</template>

<script>
import Sortable from 'sortablejs'
import APIService from '../services/APIService'

class WorkQueue {
  constructor(api) {
    this.api = api
    this.workMap = {}
  }

  enqueueUpdate(item) {
    clearTimeout(this.workMap[item.id])
    this.workMap[item.id] = setTimeout(() => {
      if (item.text !== '') {
        this.api.updateItem(item).then(item => {
          console.debug('UPDATE finished:', item)
          delete this.workMap[item.id]
        })
      }
    }, 1000)
  }
}

export default {
  name: 'List',
  data() {
    return {
      api: null,
      workQueue: null,
      focusedEl: null,
      items: []
    }
  },
  computed: {
    lastItem() {
      return this.items ? this.items[this.items.length - 1] : null
    }
  },
  methods: {
    updateDebounced(item) {
      this.workQueue.enqueueUpdate(item)
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
    addNewItem() {
      let newItem = {
        id: 'id-placeholder',
        checked: false,
        text: '',
        time_created: 'time_created-placeholder'
      }
      this.items.push(newItem)

      this.$nextTick().then(() => this.$refs[this.lastItem.id][0].focus())

      this.api
        .createItem()
        .then(item => {
          newItem.id = item.id
          newItem.time_created = item.time_created
          console.debug('CREATE finished:', newItem.id, newItem)
        })
        .catch(e => console.error('CREATE errored:', e))
    },
    _executeAsyncApiDeleteItem(item) {
      setTimeout(() => {
        if (item.id === 'id-placeholder') {
          console.debug('DELETE id-placeholder: delaying 1 second...')
          this._executeAsyncApiDeleteItem(item)
          return
        }
        this.api
          .deleteItem(item)
          .then(() => console.debug('DELETE finished:', item.id, item))
      }, 1000)
    }
  },
  mounted() {
    let el = document.getElementById('listWithHandle')
    Sortable.create(el)
    const uid = this.$route.params.id || 'default'
    console.log(`mounting List with ${uid}`)
    this.api = new APIService(uid)
    this.workQueue = new WorkQueue(this.api)
    this.api.getItems().then(items => (this.items = items))
  },
  beforeRouteUpdate(to, from, next) {
    console.debug('routing from', from, 'to', to)
    this.api = new APIService(to.params.id || 'default')
    this.workQueue = new WorkQueue(this.api)
    this.api.getItems().then(items => (this.items = items))
    next()
  }
}
</script>

<style scoped>
</style>
