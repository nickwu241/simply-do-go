<template>
  <div>
    <div>Current List: {{ lid }} </div>
    <input type="text" id="lid-input" placeholder="default" :value="uidInputDisplayValue" @keyup.enter="go">
    <button @click="go">Go</button>
    <div>
      <button @click="generateRandomId">Generate Random ID</button>
    </div>
    <div class="tooltip">
      <button @click="copyToClipboard" @mouseout="showCopiedToClipboard">
        <span class="tooltiptext" id="myTooltip">Copy to Clipboard</span>
        Share
      </button>
    </div>
  </div>
</template>

<script>
import uuidv4 from 'uuid/v4'

export default {
  name: 'Navbar',
  data() {
    return {
      lid: ''
    }
  },
  computed: {
    uidInputDisplayValue() {
      return this.lid !== 'default' ? this.lid : ''
    }
  },
  methods: {
    copyToClipboard() {
      const copyText = `https://simply-do.herokuapp.com${this.$route.path}`
      navigator.clipboard
        .writeText(copyText)
        .then(() => console.log(`copied ${copyText} to clipboard!`))
        .catch(err => console.error('error copying to clipboard:', err))

      document.getElementById('myTooltip').innerHTML =
        'Copied Link to Clipboard!'
    },
    generateRandomId() {
      document.getElementById('lid-input').value = uuidv4()
    },
    go() {
      this.lid = document.getElementById('lid-input').value || 'default'
      this.$router.push({ name: 'list', params: { id: this.lid } })
    },
    showCopiedToClipboard() {
      document.getElementById('myTooltip').innerHTML = 'Copy to Clipboard'
    }
  },
  mounted() {
    this.lid = this.$route.params.id || ''
  }
}
</script>

<style scoped>
.tooltip {
  position: relative;
  display: inline-block;
}

.tooltip .tooltiptext {
  visibility: hidden;
  width: 140px;
  background-color: #555;
  color: #fff;
  text-align: center;
  border-radius: 6px;
  padding: 5px;
  position: absolute;
  z-index: 1;
  bottom: 150%;
  left: 50%;
  margin-left: -75px;
  opacity: 0;
  transition: opacity 0.3s;
}

.tooltip .tooltiptext::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  margin-left: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: #555 transparent transparent transparent;
}

.tooltip:hover .tooltiptext {
  visibility: visible;
  opacity: 1;
}
</style>
