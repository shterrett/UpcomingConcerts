Concerts.Concert = DS.Model.extend
  title: DS.attr('string'),
  works: DS.attr('string'),
  day: DS.attr('string'),
  time: DS.attr('string'),
  link: DS.attr('string')

  safeWorks: (->
    new Handlebars.SafeString(@get('works'))
  ).property('works')
