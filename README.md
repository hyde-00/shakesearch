# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful. 

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission. 


## Submission

1. Fork this repository and send us a link to your fork after pushing your changes. 
2. Heroku hosting - The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

## Changes

1. Case insensitive searching.
2. Multiple words can be searched.
3. Display total number of matches found.
4. Display results indicating in which story/novel the query string was found.

Changes are hosted at: https://lit-journey-93419.herokuapp.com/

## Changes that could have been done with more time (priority wise)

1. Fuzzy search to accept spelling mistakes.
2. Add pagination so that results are more readable.
3. Bold the query string in the results to enable readablity.
4. On searching multi-words, even search for each word individually and add it to last pages of search.

