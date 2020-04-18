# amzn-wishlist-price-check

Check the prices of books on an Amazon public wishlist.

## Why I built this tool

### Problem

You want to know when you favorite books are on sale to get the best deal. The problem is Amazon doesn't always tell you when a book goes on sale for one or two days leaving you to miss out on a great deal.

### Solution

If you put all your books in a public whish list we can parse the list and see the prices. When we do this two to four times a day you will be able to tell when books go on sale. This makes it easy to get the best deal on your favorite books.

**NOTE:** All links and URLs use smile.amazon.com

## Requirements

1. The wishlist must be **PUBILC**

## Tunables

1. You can cover multiple wishlists and get a composite result back. We take the wishlist id and use that to fetch the list. See the example below on how to find the wish list id.
  Amazon URL: `https://smile.amazon.com/hz/wishlist/ls/3GMM56QSK63GT`
  Wishlist Id: `3GMM56QSK63GT`
2. You can set the max price you want to pay for any given book. See [Price](#price) for details

## Price

The max price comparison uses the whole dollar amount to avoid extra conversion and rounding with cents.

**Example:**
Actual book price: $2.99
Max Price value: 2

The above example will return any book that is $2.99 or lower.

## Output

The out is an HTML formatted string meant to be placed in an email.

```html
<b>Name:</b> Constitution: Book 1 of The Legacy Fleet Series
<b>Price:</b> $3.99
<a href="https://smile.amazon.com/dp/B010L6JTO0/?coliid=I17MWIZQ9RA5TN&colid=1DJLN9PNW8R59&psc=0">Buy Now</a>
```

## Timestamp

At the top of each email it will give you the time the list was last fetched

```
Books fetched on Sun Sep 15 13:23:48 CEST 2019
```
