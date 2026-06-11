# Augment
Small app as an exercise for Augment

## Requirements
1. Create a New Augment Fund
 - Each Augment Fund must have at least:
    - A name (e.g., “Augment Fund II”).
    - A total number of units (e.g., 1,000).
 - Think about how you’ll store and identify each Augment Fund (e.g., with an
integer ID, a UUID, etc.).
2. Retrieve an Augment Fund’s Current Cap Table
 - The cap table shows who owns how many units, and when they acquired those
units.
 - Each line in the cap table must include at least:
    - The owner’s name.
    - The number of units they currently own.
    - The date they acquired (or last updated) that ownership.
 - You can decide how to structure and return this information (e.g., an array of
objects, a JSON response, a CLI printout, etc.).
3. Create a New Transfer
 - A transfer is an exchange of ownership between two people for a specified
number of units within an Augment Fund.
 - You must update the Augment Fund’s cap table accordingly (the “from” person
loses some units, the “to” person gains those units).
 - You should handle relevant validations and constraints (e.g., a person can’t
transfer more units than they own).
4. Show the History of All Transfers for an Augment Fund
 - Provide a way to list all past transfers that happened in a given Augment Fund,
in some sensible order (chronological or reverse-chronological).
 - Each transfer record should, at minimum, identify which Augment Fund it pertains
to, who transferred what to whom, and when it occurred.


## How To Use This Application

### Return Format

## Improvements I Would Make For A Production Version
- Log Levels
- Robust validations (eg email)
- Email the intended recipient to invite them to Augment when someone tries to transfer shares to an Investor that does not exist.
   - Temp hold on shares for X time until the person creates account or the go back to the original owner?
   - Not sure on regulation for what's possible here
   - Probably need to confirm more info from the seller
- Need cost info on transfers
   - Pricing/most recent cost of transaction on the fund?