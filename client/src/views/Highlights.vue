<template>
  <div id="highlights">
    <div class="highlights__days-navigation">
      <div
        class="highlights__days-navigation__day"
        :class="{'highlights__days-navigation__day--selected': selectedDayOption === 'today'}"
        @click="today">Today</div>
      <div 
        class="highlights__days-navigation__day"
        :class="{'highlights__days-navigation__day--selected': selectedDayOption === 'yesterday'}"
        @click="yesterday">Yesterday</div>
      <div
        class="highlights__days-navigation__day highlights__days-navigation__custom-day"
        :class="{'highlights__days-navigation__custom-day--selected': selectedDayOption === 'custom'}">
        <v-date-picker
          v-model="day"
          :max-date="new Date()"
          :popover="{ placement: 'bottom', visibility: 'click' }">
          <div class="highlights__days-navigation__custom-day__input">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 36 36">
              <path d="M12.858 14.626h4.596v4.089h-4.596zM18.996 14.626h4.595v4.089h-4.595zM25.128 14.626h4.596v4.089h-4.596zM6.724 20.084h4.595v4.086H6.724zM12.858 20.084h4.596v4.086h-4.596zM18.996 20.084h4.595v4.086h-4.595zM25.128 20.084h4.596v4.086h-4.596zM6.724 25.54h4.595v4.086H6.724zM12.858 25.54h4.596v4.086h-4.596zM18.996 25.54h4.595v4.086h-4.595zM25.128 25.54h4.596v4.086h-4.596z"/>
              <path d="M31.974 32.198c0 .965-.785 1.75-1.75 1.75h-24c-.965 0-1.75-.785-1.75-1.75V12.099h-2.5v20.099a4.255 4.255 0 004.25 4.25h24a4.255 4.255 0 004.25-4.25V12.099h-2.5v20.099zM30.224 3.948h-1.098V2.75c0-1.517-1.197-2.75-2.67-2.75-1.474 0-2.67 1.233-2.67 2.75v1.197h-2.74V2.75c0-1.517-1.197-2.75-2.67-2.75-1.473 0-2.67 1.233-2.67 2.75v1.197h-2.74V2.75c0-1.517-1.197-2.75-2.67-2.75-1.473 0-2.67 1.233-2.67 2.75v1.197H6.224a4.255 4.255 0 00-4.25 4.25v2h32.5v-2a4.255 4.255 0 00-4.25-4.249zM11.466 7.646c0 .689-.525 1.25-1.17 1.25s-1.17-.561-1.17-1.25V2.75c0-.689.525-1.25 1.17-1.25s1.17.561 1.17 1.25v4.896zm8.08 0c0 .689-.525 1.25-1.17 1.25s-1.17-.561-1.17-1.25V2.75c0-.689.525-1.25 1.17-1.25s1.17.561 1.17 1.25v4.896zm8.08 0c0 .689-.525 1.25-1.17 1.25-.646 0-1.17-.561-1.17-1.25V2.75c0-.689.524-1.25 1.17-1.25.645 0 1.17.561 1.17 1.25v4.896z"/>
            </svg>
          </div>
        </v-date-picker>
        <div class="highlights__days-navigation__custom-day__day" v-if="customDayFormatted">{{ customDayFormatted }}</div>
      </div>
    </div>
    <div class="highlights__list">
      <highlight
        v-for="highlight in highlights"
        loading="lazy"
        :key="highlight.id"
        :id="highlight.id"
        :url="highlight.url"
        :title="highlight.title"
        :reddit-permalink="highlight.redditPermalink"
        :reddit-author="highlight.redditAuthor"
        :reddit-created-at="highlight.redditCreatedAt"
        :num-mirrors="highlight.numMirrors">
      </highlight>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from 'vuex'
import { formatDay, parseDay } from '../common.js'
import { parseISO, subDays, compareDesc, isSameDay, isFuture } from 'date-fns'
import Highlight from '../components/Highlight.vue'

const today = () => new Date()
const yesterday = () => subDays(new Date(), 1)

export default {
  components: {
    Highlight,
  },
  data () {
    return {
      day: null,
    }
  },
  methods: {
    ...mapActions(['fetchDayHighlights']),
    today () {
      this.day = today()
    },
    yesterday () {
      this.day = yesterday()
    },
    goToDay (day) {
      if (isSameDay(this.queryDay, day)) {
        return
      }
      this.$router.push({ name: 'highlights', query: { day: formatDay(day) } })
    },
    fetchHighlights () {
      this.fetchDayHighlights(this.day)
    },
    observeHighlightVideoFrames () {
      this.$nextTick(() => {
        this.frameObserver.disconnect()
        const frames = document.querySelectorAll('.highlight-video__frame')
        frames.forEach(frame => this.frameObserver.observe(frame))
      })
    },
  },
  computed: {
    ...mapState(['dayHighlights']),
    queryDay () {
      return this.$route.query.day ? parseDay(this.$route.query.day) : null
    },
    highlights () {
      if (!this.day) {
        return []
      }

      const dayFormatted = formatDay(this.day)
      if (!this.dayHighlights[dayFormatted]) {
        return []
      }

      const highlights = this.dayHighlights[dayFormatted].slice()
      highlights.sort((a, b) => compareDesc(parseISO(a.redditCreatedAt), parseISO(b.redditCreatedAt)))
      return highlights
    },
    selectedDayOption () {
      if (isSameDay(this.day, today())) {
        return 'today'
      } else if (isSameDay(this.day, yesterday())) {
        return 'yesterday'
      }
      return 'custom'
    },
    customDayFormatted () {
      return this.selectedDayOption === 'custom' ? formatDay(this.day) : null
    },
  },
  watch: {
    day () {
      this.goToDay(this.day)
    },
    queryDay () {
      this.fetchHighlights()
    },
    highlights () {
      this.observeHighlightVideoFrames()
    },
  },
  beforeCreate () {
    const onIntersectionChange = entries => entries.filter(entry => entry.isIntersecting).forEach(entry => entry.target.dispatchEvent(new Event('enteredViewport')))
    this.frameObserver = new IntersectionObserver(onIntersectionChange, { threshold: 0.0 })
  },
  created () {
    if (this.queryDay && !isFuture(this.queryDay)) {
      this.day = this.queryDay
    } else {
      this.day = today()
    }
    this.fetchHighlights()
  },
  mounted () {
    this.observeHighlightVideoFrames()
  },
}
</script>

<style lang="less">
.highlights {
  &__days-navigation {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    color: white;

    &__day {
      cursor: pointer;
      margin-right: 8px;

      &--selected {
        text-decoration: underline;
        color: black;
      }
    }

    &__custom-day {
      display: flex;
      align-items: center;

      &__day {
        margin-left: 4px;
      }

      &__input svg {
        width: 24px;
        height: 24px;
        fill: white;
      }

      &--selected {
        color: black;
        text-decoration: underline;

        svg {
          fill: black;
        }
      }
    }
  }

  &__list {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
  }
}
</style>
