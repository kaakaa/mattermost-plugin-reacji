## Emoji aliases

Some Mattermost emojis has alises. For example, `:+1:` and `:thumbsup:` indicates same emoji :+1:, and when you attach a reaction with `:thumbsup:` to a post, Mattermost considers its reaction as `:+1:`. So if you add a Reacji with `:thumbsup:`, it will never be fired.

The following table is Mattermost emoji alias list (AFAIK). You should avoid to use emojis in `aliases` row for Reacji, and should use one in `name` row.

| emoji                                                                                                    | name                            | aliases                      |
| :------------------------------------------------------------------------------------------------------- | :------------------------------ | :--------------------------- |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f606.png)                | `:laughing:`                    | `:satisfied:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f621.png)                | `:rage:`                        | `:pout:`                     |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f4a9.png)                | `:hankey:`                      | `:poop:`, `:shit:`           |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f44d.png)                | `:+1:`                          | `:thumbsup:`                 |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f44e.png)                | `:-1:`                          | `:thumbsdown:`               |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f44a.png)                | `:fist_oncoming:`               | `:facepunch:`, `:punch:`     |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/270a.png)                 | `:fist_raised:`                 | `:fist:`                     |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/270b.png)                 | `:hand:`                        | `:raised_hand:`              |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f595.png)                | `:middle_finger:`               | `:fu:`                       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f471.png)                | `:blonde_man:`                  | `:person_with_blond_hair:`   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f46e.png)                | `:policeman:`                   | `:cop:`                      |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f477.png)                | `:construction_worker_man:`     | `:construction_worker:`      |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f575-fe0f.png)           | `:male_detective:`              | `:detective:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f647.png)                | `:bowing_man:`                  | `:bow:`                      |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f481.png)                | `:tipping_hand_woman:`          | `:information_desk_person:`  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f645.png)                | `:no_good_woman:`               | `:no_good:`, `:ng_woman:`    |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f645-200d-2642-fe0f.png) | `:no_good_man:`                 | `:ng_man:`                   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f64b.png)                | `:raising_hand_woman:`          | `:raising_hand:`             |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f64e.png)                | `:pouting_woman:`               | `:person_with_pouting_face:` |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f64d.png)                | `:frowning_woman:`              | `:person_frowning:`          |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f487.png)                | `:haircut_woman:`               | `:haircut:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f486.png)                | `:massage_woman:`               | `:massage:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f46f.png)                | `:dancing_women:`               | `:dancers:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f6b6.png)                | `:walking_man:`                 | `:walking:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f3c3.png)                | `:running_man:`                 | `:runner:`, `:running:`      |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f491.png)                | `:couple_with_heart_woman_man:` | `:couple_with_heart:`        |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f46a.png)                | `:family_man_woman_boy:`        | `:family:`                   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f455.png)                | `:shirt:`                       | `:tshirt:`                   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f45e.png)                | `:mans_shoe:`                   | `:shoe:`                     |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f41d.png)                | `:bee:`                         | `:honeybee:`                 |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f42c.png)                | `:dolphin:`                     | `:flipper:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f43e.png)                | `:feet:`                        | `:paw_prints:`               |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f314.png)                | `:moon:`                        | `:waxing_gibbous_moon:`      |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f4a5.png)                | `:boom:`                        | `:collision:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f34a.png)                | `:tangerine:`                   | `:orange:`, `:mandarin:`     |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f3c4.png)                | `:surfing_man:`                 | `:surfer:`                   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f3ca.png)                | `:swimming_man:`                | `:swimmer:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f6a3.png)                | `:rowing_man:`                  | `:rowboat:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f6b4.png)                | `:biking_man:`                  | `:bicyclist:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f6b5.png)                | `:mountain_biking_man:`         | `:mountain_bicyclist:`       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f697.png)                | `:car:`                         | `:red_car:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/26f5-fe0f.png)            | `:boat:`                        | `:sailboat:`                 |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/260e-fe0f.png)            | `:phone:`                       | `:telephone:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f52a.png)                | `:hocho:`                       | `:knife:`                    |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f3ee.png)                | `:izakaya_lantern:`             | `:lantern:`                  |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/2709-fe0f.png)            | `:email:`                       | `:envelope:`                 |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f4d6.png)                | `:book:`                        | `:open_book:`                |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f4dd.png)                | `:memo:`                        | `:pencil:`                   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/2757-fe0f.png)            | `:exclamation:`                 | `:heavy_exclamation_mark:`   |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f1e8-1f1e6.png)          | `:canada:`                      | `:ca:`                       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f1ea-1f1fa.png)          | `:eu:`                          | `:european_union:`           |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f1f5-1f1f0.png)          | `:pakistan:`                    | `:pk:`                       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f1ff-1f1e6.png)          | `:south_africa:`                | `:za:`                       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/1f1ec-1f1e7.png)          | `:gb:`                          | `:uk:`                       |
| [png](https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/shipit.png)               | `:shipit:`                      | `:squirrel:`                 |
