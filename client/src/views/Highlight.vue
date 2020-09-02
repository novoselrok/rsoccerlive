<template>
  <div id="highlight">
    <highlight
      v-if="!isLoading && highlight"
      :key="highlight.id"
      :id="highlight.id"
      :url="highlight.url"
      :title="highlight.title"
      :reddit-permalink="highlight.redditPermalink"
      :reddit-author="highlight.redditAuthor"
      :reddit-created-at="highlight.redditCreatedAt">
    </highlight>
    <template v-if="mirrors.length > 0">
      <div class="highlight__mirrors-title">Mirrors</div>
      <div class="highlight__mirrors-list">
        <highlight
          v-for="mirror in mirrors"
          :key="mirror.id"
          :id="mirror.id"
          :url="mirror.url"
          :reddit-permalink="mirror.redditPermalink"
          :reddit-author="mirror.redditAuthor"
          :reddit-created-at="mirror.redditCreatedAt">
        </highlight>
      </div>
    </template>
  </div>
</template>

<script>
import { fetchApi } from '../common.js'
import Highlight from '../components/Highlight.vue'

export default {
  components: {
    Highlight,
  },
  data () {
    return {
      isLoading: false,
      highlight: null,
      mirrors: [],
    }
  },
  computed: {
    highlightId () {
      return this.$route.params.id
    },
  },
  methods: {
    fetchHighlightAndMirrors () {
      if (!this.highlightId) {
        return
      }
      this.isLoading = true
      const highlightPromise = fetchApi(`/highlights/${this.highlightId}`)
      const mirrorsPromise = fetchApi(`/highlightMirrors?highlightId=${this.highlightId}`)
      Promise.all([highlightPromise, mirrorsPromise])
        .then(([highlight, mirrors]) => {
          this.highlight = highlight
          this.mirrors = mirrors
          this.isLoading = false
        })
    },
  },
  created () {
    this.fetchHighlightAndMirrors()
  },
}
</script>

<style scoped lang="less">
.highlight {
  &__mirrors-title {
    font-size: 20px;
    margin-bottom: 16px;
  }

  &__mirrors-list {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
  }
}
</style>
