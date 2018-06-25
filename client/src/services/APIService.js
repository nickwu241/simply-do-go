export default class APIService {
  constructor(lid) {
    this.lid = lid
    this.baseUrl =
      process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : ''
  }

  getItems() {
    console.debug(`GET ${this.baseUrl}/api/list/${this.lid}/items`)
    return fetch(`${this.baseUrl}/api/list/${this.lid}/items`)
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  createItem() {
    console.debug(`POST ${this.baseUrl}/api/list/${this.lid}/items`, this.lid)
    return fetch(`${this.baseUrl}/api/list/${this.lid}/items`, {
      method: 'POST',
      body: JSON.stringify({
        checked: false,
        text: ''
      })
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  updateItem(item) {
    console.debug(`PUT ${this.baseUrl}/list/${this.lid}/api/items`, this.lid)
    return fetch(`${this.baseUrl}/list/${this.lid}/api/items` + item.id, {
      method: 'PUT',
      body: JSON.stringify(item)
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }

  deleteItem(item) {
    console.debug(
      `DELETE ${this.baseUrl}/list/${this.lid}/api/items`,
      this.lid,
      item
    )
    return fetch(`${this.baseUrl}/list/${this.lid}/api/items` + item.id, {
      method: 'DELETE',
      headers: this._defaultHeaders
    })
      .then(res => res.json())
      .catch(e => console.error(e))
  }
}
