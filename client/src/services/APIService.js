export default class APIService {
  constructor(uid) {
    this.uid = uid
    this.baseUrl =
      process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : ''
  }

  getItems() {
    console.debug(`GET ${this.baseUrl}/api/items`, this.uid)
    return fetch(`${this.baseUrl}/api/items`, {
      headers: this._defaultHeaders
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  createItem() {
    console.debug(`POST ${this.baseUrl}/api/items`, this.uid)
    return fetch(`${this.baseUrl}/api/items`, {
      method: 'POST',
      headers: this._defaultHeaders,
      body: JSON.stringify({
        checked: false,
        text: ''
      })
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  updateItem(item) {
    console.debug(`PUT ${this.baseUrl}/api/items`, this.uid)
    return fetch(`${this.baseUrl}/api/items/` + item.id, {
      method: 'PUT',
      headers: this._defaultHeaders,
      body: JSON.stringify(item)
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  deleteItem(item) {
    console.debug(`DELETE ${this.baseUrl}/api/items`, this.uid, item)
    return fetch(`${this.baseUrl}/api/items/` + item.id, {
      method: 'DELETE',
      headers: this._defaultHeaders
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  get _defaultHeaders() {
    return {
      'x-simply-do-uid': this.uid
    }
  }
}
