import Vue from 'vue'
import Vuex from 'vuex'
import VCalendar from 'v-calendar';
import { parseISO } from 'date-fns'
import { formatDay, fetchApi } from '@/common.js'

Vue.use(Vuex)
Vue.use(VCalendar)

export default new Vuex.Store({
  state: {
    dayHighlights: {},
    websocketEvents: [],
  },
  getters: {
    getNewHighlightIds (state) {
      return state.websocketEvents
        .filter(e => e.type === 'NEW_HIGHLIGHTS')
        .reduce((acc, e) => {
          return acc.concat(e.ids)
        }, [])
    },
  },
  mutations: {
    SET_DAY_HIGHLIGHTS (state, { day, highlights }) {
      Vue.set(state.dayHighlights, day, highlights)
    },
    ADD_DAY_HIGHLIGHT (state, { day, highlight }) {
      const existingDayHighlights = state.dayHighlights[day] || []
      Vue.set(state.dayHighlights, day, existingDayHighlights.concat([highlight]))
    },
    ADD_WEBSOCKET_EVENT (state, event) {
      state.websocketEvents.push(event)
    },
    RESET_WEBSOCKET_EVENTS (state) {
      state.websocketEvents = []
    },
  },
  actions: {
    fetchDayHighlights ({ state, commit }, date) {
      const dayFormatted = formatDay(date)
      if (state.dayHighlights[dayFormatted]) {
        return
      }
      fetchApi(`/highlights?day=${dayFormatted}`)
        .then(highlights => {
          commit('SET_DAY_HIGHLIGHTS', { day: dayFormatted, highlights })
        })
    },
    addWebSocketEvent ({ commit }, event) {
      commit('ADD_WEBSOCKET_EVENT', event)
    },
    resetWebSocketEvents ({ commit }) {
      commit('RESET_WEBSOCKET_EVENTS')
    },
    addNewHighlight ({ commit, state }, highlight) {
      const createdAtDate = parseISO(highlight.redditCreatedAt)
      const createdAtDayFormatted = formatDay(createdAtDate)
      const existingDayHighlightIds = (state.dayHighlights[createdAtDayFormatted] || []).map(highlight => highlight.id)
      if (!existingDayHighlightIds.includes(highlight.id)) {
        commit('ADD_DAY_HIGHLIGHT', { day: createdAtDayFormatted, highlight })
      }
    },
    addNewHighlights ({ getters, dispatch }) {
      const newHighlightsPromise = Promise.all(getters.getNewHighlightIds.map(id => fetchApi(`/highlights/${id}`)))
      newHighlightsPromise
        .then(newHighlights => {
          newHighlights.forEach(highlight => dispatch('addNewHighlight', highlight))
          dispatch('resetWebSocketEvents')
        })
    },
  },
})
