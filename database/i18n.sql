truncate i18n;



insert into
  i18n (key, de, en)
values
  ('accept_terms_and_privacy_1', 'Ich habe die', 'I have read the'),
  ('accept_terms_and_privacy_2', 'gelesen und akzeptiere diese', 'and accept them'),
  ('account_activate', 'Account aktivieren', 'Activate account'),
  ('activate', 'Aktivieren', 'Activate'),
  ('activation_code', 'Aktivierungscode', 'Activation code'),
  ('activation_code_description', 'Wir haben dir eine E-Mail mit dem Aktivierungscode geschickt. Bitte gibt diesen Code unten in das Feld ein.', 'We have sent you an e-mail with the activation code. Please enter this code in the field below.'),
  ('activation_code_send_new', 'Neuen Code senden', 'Send new code'),
  ('activation_code_send_new_description', 'Solltest du keine E-Mail bekommen haben, klicke entweder auf den Button unten um einen neuen Code anzufordern oder melde dich bei unserem', 'If you have not received an e-mail, either click on the button below to request a new code or contact our'),
  ('active', 'Aktiv', 'Active'),
  ('active_subscription_already_exists', 'Du hast bereits ein aktives Abo!', 'You already have an active subscription!'),
  ('address', 'Adresse', 'Address'),
  ('address_search', 'Adresse suchen', 'Search address'),
  ('age', 'Alter', 'Age'),
  ('alert.accept_terms_and_privacy', 'Bitte AGB und Datenschutzbedingungen lesen und akzeptieren', 'Please read and accept the terms and conditions and privacy policy'),
  ('alert.account_not_activated', 'Account noch nicht aktiviert.', 'Account not yet activated.'),
  ('alert.can_t_activate_account', 'Aktivierung des Accounts momentan nicht möglich. Bitte versuche es später nochmal.', 'Account activation not possible at the moment. Please try again later.'),
  ('alert.can_t_activeate_voucher', 'Gutschein konnte nicht eingelöst werden. Er ist entweder nicht mehr gültig, wurde schon eingelöst oder existiert nicht.', 'Voucher could not be redeemed. It is either no longer valid, has already been redeemed or does not exist.'),
  ('alert.can_t_change_email', 'Momentan kannst du deine E-Mail nicht ändern. Bitte versuche es später nochmal.', 'You cannot change your email at the moment. Please try again later.'),
  ('alert.can_t_change_industry', 'Deine Branche kann moment nicht geändert werden, bitte versuche es später nochmal.', 'Your industry cannot be changed at the moment, please try again later.'),
  ('alert.can_t_change_password', 'Das Passwort kann momentan nicht geändert werden. Bitte versuche es später nochmal.', 'The password cannot be changed at the moment. Please try again later.'),
  ('alert.can_t_change_phone', 'Deine Telefonnummer kann moment nicht geändert werden, bitte versuche es später nochmal.', 'Your phone number cannot be changed at the moment, please try again later.'),
  ('alert.can_t_change_rating', 'Kann Bewertung momentan nicht ändern, bitte später noch einmal versuchen.', 'Cannot change rating at the moment, please try again later.'),
  ('alert.can_t_change_username', 'Dein Benutzername kann moment nicht geändert werden, bitte versuche es später nochmal.', 'Your username cannot be changed at the moment, please try again later.'),
  ('alert.can_t_chante_taxid', 'Deine Steuernummer ID kann moment nicht geändert werden, bitte versuche es später nochmal.', 'Your tax number ID cannot be changed at the moment, please try again later.'),
  ('alert.can_t_conclude_subscription', 'Momentan können keine Abonnements abgeschlossen werden. Bitte versuche es später nochmal.', 'No subscriptions can be taken out at the moment. Please try again later.'),
  ('alert.can_t_create_account', 'Leider können momentan keine neuen Accounts angelegt werden, bitte versuche es später noch einmal.', 'Unfortunately no new accounts can be created at the moment, please try again later.'),
  ('alert.can_t_create_deal', 'Es können momentan keine Deals erstellt werden. Wir arbeiten bereits mit Hochdruck an einer Lösung. Bitte versuche es später noch einmal.', 'No deals can be created at the moment. We are already working hard on a solution. Please try again later.'),
  ('alert.can_t_delet_profile_image', 'Profilbild kann nicht gelöscht werden, bitte versuche es später noch einmal.', 'Profile picture cannot be deleted, please try again later.'),
  ('alert.can_t_delete_image', 'Konnte Bild nicht löschen, bitte später nochmal versuchen.', 'Could not delete image, please try again later.'),
  ('alert.can_t_delete_rating', 'Konnte Bewertung nicht löschen, bitte später noch einmal versuchen.', 'Could not delete rating, please try again later.'),
  ('alert.can_t_load_address', 'Kann aktuelle Adresse momentan nicht laden, bitte versuche es später nochmal.', 'Cannot load current address at the moment, please try again later.'),
  ('alert.can_t_load_deal', 'Der Deal konnte leider nicht gefunden werden, bitte versuche es später nochmal.', 'Sorry, the deal could not be found, please try again later.'),
  ('alert.can_t_load_deal_images', 'Kann Deal Bilder momentan nicht laden, bitte versuche es später nochmal.', 'Cannot load deal images at the moment, please try again later.'),
  ('alert.can_t_load_dealer_images', 'Kann Dealer Bilder momentan nicht laden, bitte versuche es später nochmal.', 'Cannot load dealer images at the moment, please try again later.'),
  ('alert.can_t_load_favorite_deals', 'Kann favorisierte Dealer Deals nicht laden, bitte später nochmal versuchen.', 'Cannot load favorite dealer deals, please try again later.'),
  ('alert.can_t_load_filter', 'Leider können die Filter gerade nicht geladen werden, bitte versuche es später noch einmal.', 'Unfortunately the filters cannot be loaded at the moment, please try again later.'),
  ('alert.can_t_load_profile_image', 'Kann Profilbild nicht laden, bitte versuche es später nochmal.', 'Cannot load profile picture, please try again later.'),
  ('alert.can_t_load_settings', 'Kann Einstellungen gerade nicht laden, bitte versuche es später noch einmal.', 'Cannot load settings at the moment, please try again later.'),
  ('alert.can_t_load_statistics', 'Momentan können keine Statistiken abgerufen werden, bitte versuche es später noch einmal.', 'Currently no statistics can be retrieved, please try again later.'),
  ('alert.can_t_load_top_deals', 'Kann momentan die Top Deals nicht laden, bitte versuche es später nochmal.', 'Cannot load the Top Deals at the moment, please try again later.'),
  ('alert.can_t_rate', 'Momentan können keine Bewertungen abgegeben werden, bitte versuche es später noch einmal.', 'No reviews can be submitted at the moment, please try again later.'),
  ('alert.can_t_report_deal', 'Nur angemeldete User können einen Deal melden', 'Only registered users can report a deal'),
  ('alert.can_t_save_deal', 'Leider ist beim Erstellen etwas schief gegangen, bitte versuche es später nochmal.', 'Unfortunately something went wrong while creating the deal, please try again later.'),
  ('alert.can_t_save_fav_dealer', 'Kann favorisierte Dealer momentan nicht speichern, bitte versuche es später nochmal.', 'Cannot save favorite dealers at the moment, please try again later.'),
  ('alert.can_t_save_filter', 'Deine Filtereinstellungen konnten nicht speichern werden, bitte versuche es später noch einmal.', 'Your filter settings could not be saved, please try again later.'),
  ('alert.can_t_save_message', 'Leider ist beim speichern deiner Nachricht etwas schief gegangen, bitte versuche es später noch einmal.', 'Unfortunately something went wrong while saving your message, please try again later.'),
  ('alert.can_t_save_profile_image', 'Profilbild kann nicht gespeichert werden, bitte versuche es später noch einmal.', 'Profile picture cannot be saved, please try again later.'),
  ('alert.can_t_save_rating', 'Konnte Bewertung nicht speichern, bitte versuche es später noch mal.', 'Could not save rating, please try again later.'),
  ('alert.can_t_save_report', 'Der Deal konnte leider nicht gemeldet werden, bitte versuche es später noch einmal.', 'Sorry, the deal could not be reported, please try again later.'),
  ('alert.can_t_save_settings', 'Kann Einstellung leider nicht speichern, bitte versuche es später nochmal.', 'Cannot save setting, please try again later.'),
  ('alert.can_t_send_activation_email', 'Momentan können keine Aktivierungs-Emails versendet werden. Bitte versuche es später nochmal.', 'No activation emails can be sent at the moment. Please try again later.'),
  ('alert.can_t_update_address', 'Kann Adresse momentan nicht aktuallisieren, bitte versuche es später noch einmal.', 'Cannot update address at the moment, please try again later.'),
  ('alert.email_already_used', 'E-Mail ist bereits mit einem anderen Account verknüpft.', 'Email is already linked to another account.'),
  ('alert.enter_description', 'Bitte eine Beschreibung angeben', 'Please enter a description'),
  ('alert.enter_title', 'Bitte einen Titel angeben', 'Please enter a title'),
  ('alert.error_while_signup', 'Fehler bei der Registrierung, bitte versuche es später nochmal.', 'Error during signup, please try again later.'),
  ('alert.invalid_activation_code', 'Der angegebene Aktivierungscode ist ungütlig.', 'The activation code is invalid.'),
  ('alert.invalid_address', 'Ungültige Adresse, bitte geben Sie eine genauere Adresse an.', 'Invalid address, please enter a more accurate address.'),
  ('alert.invalid_email', 'Bitte eine gültige E-Mail angeben', 'Please enter a valid e-mail address'),
  ('alert.invalid_username_or_password', 'Benutzername oder Passwort falsch', 'Username or password incorrect'),
  ('alert.login_not_possible', 'Login gerade nicht möglich, bitte später nochmal versuchen.', 'Login currently not possible, please try again later.'),
  ('alert.message_delay', 'Du kannst uns nur alle 5 Minuten eine neue Nachricht schreiben, bitte versuche es später noch einmal.', 'You can only send us a new message every 5 minutes, please try again later.'),
  ('alert.not_your_statistics', 'Sie sind nicht berechtigt die Statistiken dieses Deals zu sehen!', 'You are not authorized to see the statistics of this deal!'),
  ('alert.password_repeat_not_matching', 'Das Passwort und die Wiederholung stimmen nicht überein', 'The password and the repetition do not match'),
  ('alert.provide_email', 'Bitte eine gültige E-Mail angeben', 'Please enter a valid e-mail address'),
  ('alert.provide_runtime_or_enddate', 'Bitte entweder eine Laufzeit oder ein Enddatum angeben', 'Please provide either a runtime or an end date'),
  ('alert.provide_start_date', 'Bitte ein (gültiges) Startdatum angeben', 'Please provide a (valid) start date'),
  ('alert.provide_username', 'Bitte einen Benutzernamen angeben', 'Please enter a username'),
  ('alert.report_cause', 'Bitte gib an, was an dem Deal nicht passt', 'Please specify what is wrong with the deal'),
  ('alert.select_category', 'Bitte eine Kategorie auswählen', 'Please select a category'),
  ('alert.start_before_end', 'Das Startdatum muss vor dem Enddatum liegen', 'The start date must be before the end date'),
  ('alert.technical_problem', 'Leider gibt es aktuell ein technisches Problem, bitte versuche es später noch einmal.', 'Unfortunately there is currently a technical problem, please try again later.'),
  ('alert.username_already_exists', 'Der Benutzername ist leider schon vergeben', 'Unfortunately the username is already taken'),
  ('alert.username_or_email_already_used', 'Benutzername oder E-Mail bereits vergeben', 'Username or e-mail already taken'),
  ('all', 'Alle', 'All'),
  ('already_have_active_subscription', 'Du hast bereits ein aktives Abo!', 'You already have an active subscription!'),
  ('and', 'und', 'and'),
  ('back', 'Zurück', 'Back'),
  ('basic', 'BASIC', 'BASIC'),
  ('basic_pricing', 'Basic', 'Basic'),
  ('basic_pricing_description_1', 'Zahle pro Tagesdeal', 'Pay per day'),
  ('basic_pricing_description_1_24_hours', '(24 Std. Laufzeit)', '(24h runtime)'),
  ('basic_pricing_description_2', 'Kein Abo / keine Grundgebühr', 'No subscription / no basic fee'),
  ('basic_vs_pro', 'Basic vs Pro', 'Basic vs Pro'),
  ('cancel', 'Abbrechen', 'Cancel'),
  ('category', 'Kategorie', 'Category'),
  ('category_select', 'Kategorie auswählen', 'Select category'),
  ('change', 'Ändern', 'Change'),
  ('city', 'Ort', 'City'),
  ('common', 'Allgemein', 'General'),
  ('contact_us', 'Kontaktiere uns', 'Contact us'),
  ('contact_us_description', 'Hast du Fragen, Anmerkungen oder Lob? Oder willst du sonst etwas los werden? Schreib uns gerne eine Narchricht!', 'Do you have any questions, comments or praise? Or do you have something else to say? Feel free to write us a message!'),
  ('contact_us_restriction', 'Auch wenn wir uns riesig über jedes Feedback freuen, kannst du uns nur alle 5 Minuten eine neue Nachricht schicken.', 'Even though we are always delighted to receive feedback, you can only send us a new message every 5 minutes.'),
  ('continue_as_basic_user', 'Weiter als Basic User', 'Continue as Basic User'),
  ('costs', 'Kosten', 'Costs'),
  ('create', 'Erstellen', 'Create'),
  ('days', 'Tag(e)', 'Day(s)'),
  ('deal', 'Deal', 'Deal'),
  ('deal_already_reported', 'Du hast diesen Deal schon gemeldet', 'You have already reported this deal'),
  ('deal_create', 'Deal erstellen', 'Create deal'),
  ('deal_rework', 'Nochmal überarbeiten', 'Edit deal'),
  ('deal_start_description', 'Dein Deal startet sofort, wenn du auf "Erstellen" klickst!', 'Your deal starts immediately when you click on "Create"!'),
  ('deal_start_now', 'Jetzt starten!', 'Start now!'),
  ('deal_summary', 'Das Wichtigste zusammengefasst', 'The most important things summarized'),
  ('deal_summary_active_voucher', 'Du sparst durch einen aktiven Gutschein', 'You save with an active voucher'),
  ('deal_summary_error_1', 'Wir können momentan nicht überprüfen, ob du ein aktives Abo oder aktiven Gutschein hast.', 'We are currently unable to check whether you have an active subscription or active voucher.'),
  ('deal_summary_error_2', 'Deshalb müssen wir dir für diesen Deal leider den Standard-Preis berechnen.', 'That''s why we unfortunately have to charge you the standard price for this deal.'),
  ('deal_summary_error_3', 'Du kannst es auch später nochmal versuchen oder melde dich bitte bei unserem Support-Team:', 'You can also try again later or please contact our support team:'),
  ('deal_summary_error_4', 'Nachricht schreiben', 'Write us a message'),
  ('deal_summary_no_free_days_left', 'Du hast deine freien Tage für diese Monatsperiode aufgebraucht. Falls du mehr freie Tage benötigst, kannst du auf ein größeres Abo umsteigen', 'You have used up your free days for this monthly period. If you need more days off, you can upgrade to a larger subscription'),
  ('dealer_images_add', 'Füge ein paar Bilder hinzu und mach deine Seite noch schöner! 🚀', 'Add a few pictures and make your page even more beautiful! 🚀'),
  ('dealer_images_no_images_yet', 'Leider gibt es hier noch keine Bilder vom Dealer zu bestaunen.', 'Unfortunately there are no pictures of the dealer to marvel at yet.'),
  ('deals_left', 'Deals verbleibend', 'deals left'),
  ('delete', 'Löschen', 'Delete'),
  ('description', 'Beschreibung', 'Description'),
  ('deselect_all', 'Alles abwählen', 'Deselect all'),
  ('duration', 'Dauer', 'Duration'),
  ('email', 'E-Mail', 'E-Mail'),
  ('email_change', 'E-Mail ändern', 'Change e-mail'),
  ('email_change_description', 'Bitte gibt die neue E-Mail Adresse ein. Wir schicken an diese Adresse eine E-Mail mit einem Bestätigungslink zum aktivieren der neuen E-Mail Adresse.', 'Please enter the new e-mail address. We will send an e-mail to this address with a confirmation link to activate the new e-mail address.'),
  ('end', 'Ende', 'End'),
  ('ends_on', 'Endet am', 'Ends on'),
  ('exclusive', 'Exclusive', 'Exclusive'),
  ('filter', 'Filter', 'Filter'),
  ('free_days_left', 'Verbleibende freie Tage', 'Free days left'),
  ('from', 'Von', 'From'),
  ('gender', 'Geschlecht', 'Gender'),
  ('house_number', 'Hausnummer', 'House number'),
  ('i_am_a_dealer', 'Ich bin ein Dealer', 'I am a dealer'),
  ('images_add', 'Bilder hinzufügen', 'Add images'),
  ('imprint', 'Impressum', 'Imprint'),
  ('industry', 'Branche', 'Industry'),
  ('industry_select', 'Branche auswählen', 'Select industry'),
  ('language', 'Sprache', 'Language'),
  ('location_change_description', 'Bitte prüfe genau, ob die Adresse auf der Karte korrekt angezeigt wird. Es ist extrem wichtig, dass die Position auf der Karte stimmt, da hier die eingestellten Deals angezeigt werden!', 'Please check carefully that the address is displayed correctly on the map. It is extremely important that the position on the map is correct, as this is where your deals are displayed!'),
  ('login', 'Anmelden', 'Login'),
  ('logout', 'Abmelden', 'Logout'),
  ('manage_subscription', 'Abo verwalten', 'Manage subscription'),
  ('message', 'Nachricht', 'Message'),
  ('month', 'Monat', 'Month'),
  ('monthly_exclusive_1', '90 kostenlose Tagesdeals pro Monat', '90 free daily deals per month'),
  ('monthly_exclusive_2', 'Preisvorteil von ~44% je Deal', 'Price advantage of ~44% per deal'),
  ('monthly_exclusive_3', 'Einfache Statistiken', 'Simple statistics'),
  ('monthly_exclusive_4', 'Monatlich kündbar', 'Monthly cancellable'),
  ('monthly_premium_1', '300 kostenlose Tagesdeals pro Monat', '300 free daily deals per month'),
  ('monthly_premium_2', 'Preisvorteil von ~76% je Deal', 'Price advantage of ~76% per deal'),
  ('monthly_premium_3', 'Erweiterte Statistiken (bald verfügbar)', 'Advanced statistics (available soon)'),
  ('monthly_premium_4', 'Monatlich kündbar', 'Monthly cancellable'),
  ('monthly_starter_1', '30 kostenlose Tagesdeals pro Monat', '30 free daily deals per month'),
  ('monthly_starter_2', 'Preisvorteil von ~33% je Deal', 'Price advantage of ~33% per deal'),
  ('monthly_starter_3', 'Monatlich kündbar', 'Monthly cancellable'),
  ('monthly_subscription', 'Monatsabo', 'Monthly subscription'),
  ('ok', 'OK', 'OK'),
  ('optional', 'optional', 'optional'),
  ('password', 'Passwort', 'Password'),
  ('password_change', 'Passwort ändern', 'Change password'),
  ('password_change_description', 'Bitte gibt die E-Mail-Adresse ein, an die wir den Code zum Ändern deines Passwords senden sollen.', 'Please enter the e-mail address to which we should send the code to change your password.'),
  ('password_new', 'Neues Passwort', 'New password'),
  ('password_repeat', 'Passwort wiederholen', 'Repeat password'),
  ('password_reset', 'Passwort zurücksetzen', 'Reset password'),
  ('past', 'Abgelaufen', 'Past'),
  ('per_deal_day', 'pro Tag', 'per day'),
  ('perimeter', 'Umkreis', 'Perimeter'),
  ('phone', 'Telefon', 'Phone'),
  ('planed', 'Geplant', 'Planed'),
  ('premium', 'Premium', 'Premium'),
  ('prices_and_subscriptions', 'Preise & Abos', 'Prices & Subscriptions'),
  ('pricing', 'Preise & Abos', 'Pricing'),
  ('privacy', 'Datenschutz', 'Privacy'),
  ('pro', 'PRO', 'PRO'),
  ('profile_picture', 'Profilbild', 'Profile picture'),
  ('rate', 'Bewerten', 'Rate'),
  ('rating', 'Bewertung', 'Rating'),
  ('rating_add', 'Bewertung abgeben', 'Submit rating'),
  ('rating_be_the_first_1', 'Sei der erste der eine Bewertung schreibt!', 'Be the first to leave a rating!'),
  ('rating_be_the_first_2', '(nur für PRO-Mitglieder)', '(only for PRO members)'),
  ('rating_not_yet_rated', 'Leider hast du noch keine Bewertung bekommen 😕', 'Unfortunately you haven not received a rating yet 😕'),
  ('rating_text', 'Bewertungstext', 'Rating text'),
  ('redeem', 'Einlösen', 'Redeem'),
  ('report', 'Melden', 'Report'),
  ('runtime', 'Laufzeit', 'Runtime'),
  ('runtime_individual', 'Individuelle Laufzeit', 'Individual runtime'),
  ('save', 'Speichern', 'Save'),
  ('save_10_percent_with_yearly', 'Spare weitere 10% beim Abschluss eines Jahresabos!', 'Save a further 10% when you choose a yearly subscription!'),
  ('save_additionally_as_template', 'Zusätzlich als Vorlage speichern', 'Additionally save as template'),
  ('save_money', 'Spare Geld', 'Save money'),
  ('select_all', 'Alles auswählen', 'Select all'),
  ('send', 'Senden', 'Send'),
  ('settings', 'Einstellungen', 'Settings'),
  ('show_on_map', 'Auf Karte anzeigen', 'Show on map'),
  ('signup', 'Registrieren', 'Signup'),
  ('signup_complete', 'Registrierung abschließen', 'Complete signup'),
  ('signup_now_1', 'Registriere dich jetzt', 'Signup now'),
  ('signup_now_2', 'um Dealer bewerten zu können!', 'to be able to rate dealers!'),
  ('start', 'Start', 'Start'),
  ('start_immediately', 'Sofort starten', 'Start immediately'),
  ('starter', 'Starter', 'Starter'),
  ('street', 'Straße', 'Street'),
  ('subscribe', 'schließe ein Abo ab!', 'subscribe!'),
  ('subscription_buy_now', 'Jetzt Abo abschließen und Geld sparen!', 'Buy a subscription now and save money!'),
  ('subscription_no_info_available', 'Kann aktuellen Status des Abos momentan nicht abrufen', 'Cannot retrieve current subscription status at the moment'),
  ('subscription_not_active', 'Kein Abo abgeschlossen', 'No subscription'),
  ('subscription_subscribe', 'Abo abschließen', 'Subscribe'),
  ('tax_id', 'Umsatzsteuer ID', 'Tax ID'),
  ('tell_us_what_is_wrong', 'Sag uns, was an dem Deal nicht passt', 'Tell us what is wrong with this deal'),
  ('templates', 'Vorlagen', 'Templates'),
  ('terms', 'AGB', 'Terms'),
  ('title', 'Titel', 'Title'),
  ('top_25', 'Top 25', 'Top 25'),
  ('top_deals_in_your_area', 'TOP-Deals in deiner Nähe', 'TOP deals in your area'),
  ('until', 'Bis', 'Until'),
  ('use_current_location', 'Aktuellen Standort verwenden', 'Use current location'),
  ('username', 'Benutzername', 'Username'),
  ('voucher_already_activated', 'Bereits eingelöste Gutscheine', 'Already redeemed vouchers'),
  ('voucher_code', 'Gutscheincode', 'Voucher code'),
  ('voucher_error_cannot_show', 'Kann aktive Gutscheine momentan nicht anzeigen, bitte versuche es später nochmal.', 'Cannot display active vouchers at the moment, please try again later.'),
  ('vouchers', 'Gutscheine', 'Voucher'),
  ('year', 'Jahr', 'Year'),
  ('yearly_exclusive_1', '90 kostenlose Tagesdeals pro Monat', '90 free daily deals per month'),
  ('yearly_exclusive_2', 'Preisvorteil von ~50% je Deal', 'Price advantage of ~50% per deal'),
  ('yearly_exclusive_3', 'Einfache Statistiken', 'Simple statistics'),
  ('yearly_exclusive_4', 'Monatlich kündbar', 'Cancellable at the end of the subscription'),
  ('yearly_premium_1', '300 kostenlose Tagesdeals pro Monat', '300 free daily deals per month'),
  ('yearly_premium_2', 'Preisvorteil von ~79% je Deal', 'Price advantage of ~79% per deal'),
  ('yearly_premium_3', 'Erweiterte Statistiken (bald verfügbar)', 'Advanced statistics (available soon)'),
  ('yearly_premium_4', 'Monatlich kündbar', 'Cancellable at the end of the subscription'),
  ('yearly_starter_1', '30 kostenlose Tagesdeals pro Monat', '30 free daily deals per month'),
  ('yearly_starter_2', 'Preisvorteil von ~39% je Deal', 'Price advantage of ~39% per deal'),
  ('yearly_starter_3', 'Kündbar zum Aboende', 'Cancellable at the end of the subscription'),
  ('yearly_subscription', 'Jahresabo', 'Yearly subscription'),
  ('you_want_to_report', 'Du willst den folgenden Deal melden?', 'You want to report the following deal?'),
  ('your_location', 'Dein Standort', 'Your location'),
  ('zipcode', 'PLZ', 'ZIP Code')
