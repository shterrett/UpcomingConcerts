Ember.Router.map ->
  @resource 'concerts',
    path: '/'

Concerts.ConcertsRoute = Ember.Route.extend
  model: ->
    console.log 'route model'
    @store.find('concert')
