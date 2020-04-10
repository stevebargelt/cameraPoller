# Litterbox Camera Poller

## Overview

This project attempts to determine which of our cats is using the litterbox, if they were headed in our out of the box and record the event so that the cat's strange cat parents can keep track of litterbox habits.

How you might ask? The app attempts to capture a still photo of a cat headed into or out of the litterbox. It does this by snapping a photo when the camera detects motion. It then sends the image to an Azure Custom Vision service for classification.

Once the image is tagged with predictions the app attempts to decide what just happened. Was it a false motion event (Negative)? Flipping lights switches on and off cause motion events with no cat. Our litterbox is in the laundry room, so we humans use the room, which also cause (Negative) false motion events.

If there was a cat in the picture according to the Custom Vision Service, next we attempt to determine if the cat was headed into or out of he litterbox with another Azure Custom Vision Project. Once those results are in we send the best picture of the bunch (highest probability according to the Custom Vision service) to Firebase Storage, then write the rest of the details to a Firebase Firestore database.

Ultimately the cat parents can see in an iOS app which cat is using the litterbox and when. It's over the top, I know.

## Slightly more detail

This particular project:

1. Polls a local IP camera for motion once every n seconds
1. If there is motion it saves an image locally
1. If a pic is saved it sends it off to the Azure Custom Vision service for a Cat prediction
1. If there was a cat in the pic, it tries to determine if the cat was headed into our out of the litterbox
1. Once an image set is collected (arbitrary number probably between 3-5) or a timeout is reached:
1. It picks the highest probability cat image from the bunch
1. Saves the image to Firebase Storage
1. Saves the "Trip" details to Firestore
1. Goes back to polling the camera

## Disclaimer

This system works for us and our situation. It is pretty neat but probably not useful to anyone else. Unless it sparks your creativity and curiosity, then it was plenty useful.

## Why

Cat Parents, that's why. Well that and an incurable Geek who is a Cat Dad.

This started when we decided that our cats needed two litterboxes instead of one. This decision was based purely on our laziness; we wanted to change the litter less often. Once we had two litterboxes in place I became curious, as one does, whether the cats each used their own litterbox of if they alternated between the two. I had the brilliant idea to repurpose a "security" camera so I could determine if there was a pattern. Then our laziness vaulted to a new level when we purchased an automatic scooping litterbox.

Our LitterRobot replaced the two litterboxes before my little experiment took off. It seemed like an extravagant purchase but I will never go back to scooping, or let's face it tossing 20 pounds of litter a week because we were too lazy (busy! we are busy!) to scoop daily. I thought my geek dream was dead. One of the cool things about the LitterRobot is that, through the mobile app you can see how often your cats use the litterbox. The downside is that it can't tell **which** cat used the box so my dream was, once again, alive.

I wanted to know which cat was using the litterbox when I received that push notification from the LitterRobot app, duh!

I suppose I could make up some reasons like by seeing pattern changes in litter usage one could tell if a cat was sick. That is an excuse I came up with later. Although, I suppose this data could be used for that.
