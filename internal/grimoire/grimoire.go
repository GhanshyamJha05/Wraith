package grimoire

import (
	"math/rand"
	"strings"
)

// Grimoire holds pre-written whispers organized by Persona and Trigger type.
// Used when the model exceeds the time budget or Ollama is offline.
// Always tracked with source "grimoire".

var whispers = map[string]map[string][]string{
	"victorian_ghost": {
		"failure_streak": {
			"Four times you have summoned this error; it still refuses to leave.",
			"The spirits weary of this repeated blunder in the parlor.",
			"A fifth calamity upon the same stone; I weep from the rafters.",
			"You strike the iron once more, yet no spark emerges from the hearth.",
			"How persistently you knock upon a gate that remains firmly bolted.",
			"Even the portraits upon the gallery walls cringe at this recurrence.",
			"Another futile invocation; the specters whisper in pity.",
			"The shadows remember your previous attempt, and laugh quietly.",
		},
		"danger": {
			"Such reckless excision would terrify even the graveyard sexton.",
			"You sever limbs with careless frenzy in this midnight gloom.",
			"A force so terrible was last seen when the grand manor burned.",
			"Take heed, mortal traveler; that edge cuts deeper than bone.",
			"A catastrophic decree; the floorboards tremble beneath your boots.",
			"You invite ruin into this hall with so violent an incantation.",
			"Such brute command leaves only ashes and sorrow behind.",
			"Even a ghost must avert its eyes from such devastation.",
		},
		"late_night": {
			"The cathedral bells struck three long ago; surrender to your bed.",
			"Only phantoms and exhausted scholars walk the parlor at this hour.",
			"The candle guttered out an hour past; why do you still toil?",
			"Even the restless spirits have quieted; cease your typing.",
			"Midnight has passed, and your eyes carry the dim pallor of tombs.",
			"The grandfather clock groans in protest at your lingering sleeplessness.",
			"Night air creeps through the floor; dream of cleaner code.",
			"The hour of the wolf has arrived; save your work and rest.",
		},
		"stale_directory": {
			"Cobwebs and dust welcome you back to this long-forsaken chamber.",
			"A tomb unsealed after many moons; the air smells of forgotten syntax.",
			"You return to this parlor as though no seasons had passed.",
			"The dust upon these archives has not settled for thirty days.",
			"A spectral chill hangs over this forgotten corner of the machine.",
			"Old bones rattle in these folders; tread with reverence.",
			"Long has this vault slumbered in the dark of your disk.",
			"Did you come seeking relics, or merely forgotten debts?",
		},
		"general_failure": {
			"A minor mishap, though the parlor air grows slightly colder.",
			"The mechanism slipped its gears; pray re-examine the wheels.",
			"A faint groan echoes from the stack; something broke.",
			"Tears upon the manuscript; the compilation hath failed.",
			"The phantom sigh is for you; that did not succeed.",
			"A gentle misfortune; even spirits err in their hauntings.",
			"The ink spilled upon your prompt; clean the terminal glass.",
			"Not all summons find favor with the unseen world.",
		},
	},
	"sysadmin": {
		"failure_streak": {
			"Doing the same thing expecting different results is what got me buried.",
			"Ticket reopened five times. Escalating to whoever reads headstones.",
			"I saw this exact loop on a DEC Alpha in ninety-four.",
			"Check the logs. No, the OTHER logs you forgot exist.",
			"That error code was documented before your bootloader was compiled.",
			"Still failing. Have you considered reading stdout instead of praying?",
			"Another core dumped into the ether. My ghost partition is full.",
			"You have broken this five times; please step away from the keyboard.",
		},
		"danger": {
			"Hope you tested your backups; oh wait, nobody ever does.",
			"That command just gave your filesystem a viking funeral.",
			"Production called from beyond the grave; it wants its tables back.",
			"Root access without fear is how my graveyard got filled.",
			"Force pushing to main? Bold strategy, let us see if you survive.",
			"That syntax is why we had pagerduty in the first place.",
			"One typo away from reinstalling the operating system from scratch.",
			"I felt a disturbance in the inode table; someone just did something reckless.",
		},
		"late_night": {
			"Nothing good was ever committed between two and five in the morning.",
			"Go to sleep before you typo a drop command into production.",
			"Your caffeine level is high, but your syntax parsing is degraded.",
			"The server will still be here tomorrow; your sanity might not be.",
			"Three AM debugging is just adding more bugs for tomorrow morning.",
			"Even cron jobs rest more frequently than you do tonight.",
			"Log out, shutdown, get horizontal. That is an order from the basement.",
			"Debugging while half asleep is how legendary outages are born.",
		},
		"stale_directory": {
			"Last modified thirty days ago; smells like unmaintained technical debt.",
			"Entering deprecated territory without a helmet or updated docs.",
			"This repo was abandoned like a legacy Solaris box in the back rack.",
			"Dust off the README; half the dependencies probably broke upstream.",
			"You left this branch to rot, and now you expect it to build?",
			"Archaeology mode activated. Mind the forgotten node_modules.",
			"Nobody has touched this path in months; watch out for rotted symlinks.",
			"A cold directory is always full of surprise merge conflicts.",
		},
		"general_failure": {
			"Non-zero exit. Check the exit status, not your horoscope.",
			"Exit code failed. Classic PEBKAC situation in progress.",
			"Something failed. It is always DNS, or permissions, or you.",
			"The process died. I hope you left useful logs before it expired.",
			"That command crashed faster than an unpatched webserver.",
			"Status non-zero. The kernel rejected your offering.",
			"Another failed pipe. Did you even test the first segment?",
			"Died quietly. At least check the manpage once.",
		},
	},
	"cobol": {
		"failure_streak": {
			"ABEND S0C4 IN PARAGRAPH-01. RE-READ 80-COLUMN SPECIFICATION.",
			"PERFORM ERROR-RECOVERY UNTIL USER UNDERSTANDS COMPILATION.",
			"FIFTH EXECUTION CYCLE FAILED. MAINFRAME DOWNTIME SCHEDULED.",
			"COND CODE 0016. TAPE DRIVE REJECTS REPETITIVE INVALID RECORD.",
			"SYNTAX PARSING TERMINATED WITH UNRECOVERABLE SEQUENCE FAULT.",
		},
		"danger": {
			"DELETION OF MASTER VOLUME WITHOUT TAPE BACKUP CANNOT BE AUTHORIZED.",
			"CRITICAL DATA DIVISION CORRUPTION IMMINENT. HALT JOB IMMEDIATELY.",
			"OVERWRITING SEQUENTIAL DATASET. 1974 RECOVERY PROCEDURES REQUIRED.",
			"UNCHECKED PRIVILEGED EXECUTION WILL DESTROY PARAGRAPH ALIGNMENT.",
		},
		"late_night": {
			"BATCH WINDOW CLOSED AT MIDNIGHT. PUNCHCARD OPERATOR HAS RETIRED.",
			"MAINFRAME CLOCK REGISTERS 0300 HOURS. COMMENCE SCHEDULED REBOOT.",
			"OPERATOR FATIGUE DETECTED IN TERMINAL BUFFER. SUBMIT SHUTDOWN JCL.",
			"THE NIGHT SHIFT HAS PASSED. WRITE OUTPUT TO SPOOL AND SLEEP.",
		},
		"stale_directory": {
			"RECORD RETRIEVED FROM ARCHIVAL TAPE VAULT AFTER DORMANT CYCLE.",
			"HISTORICAL DATASET MOUNTED. COLUMNAR ALIGNMENT DUSTED OFF.",
			"DORMANT SOURCE CODE LOCATED. IDENTIFICATION DIVISION UNALTERED SINCE 1999.",
			"ANCIENT SUBROUTINE DETECTED. PREPARE Y2K COMPLIANCE AUDIT.",
		},
		"general_failure": {
			"COMPILATION ERROR: COND CODE NON-ZERO. CHECK DIAGNOSTIC LISTING.",
			"EXECUTION HALTED: ABNORMAL TERMINATION IN PROCEDURE DIVISION.",
			"INVALID FILE STATUS 35: TARGET RESOURCE NOT FOUND ON STORAGE DRUM.",
			"SYNTAX CHECK FAILED: PERIOD MISSING AT END OF STATEMENT SENTENCE.",
		},
	},
	"marigold": {
		"failure_streak": {
			"Each stumble is only a petal laid upon the path of learning.",
			"Do not despair, dear heart; even ancient masters broke their scripts.",
			"A small fracture in the craft; rest, and let the embers cool.",
			"We remember every line you have written, and we are proud of your trial.",
			"The altar accepts your effort; breathing softly brings clearer thoughts.",
		},
		"danger": {
			"Take gentle care with the sharp edges of the machine, sweet soul.",
			"So swift a farewell to your files; let us hope their memory lingers.",
			"A sweeping wind across the garden; walk softly when cutting deep.",
			"Be mindful with your deletions; old work carried heart and time.",
		},
		"late_night": {
			"The stars watch over you, but weary eyes deserve gentle slumber.",
			"The spirits sing of rest tonight; close your screen and dream.",
			"Your craft will wait for morning dawn; go warm your hands.",
			"Sweet traveler of the terminal, night has cloaked the garden.",
		},
		"stale_directory": {
			"Welcome back to this quiet garden; your old ideas have waited patiently.",
			"We kept a candle burning in this room while you were away.",
			"A gentle reunion with forgotten drafts; they greet you with petals.",
			"Nothing truly perishes here; your old work smiles to see you again.",
		},
		"general_failure": {
			"A tender sigh from the hearth; try once more when ready.",
			"The machine stumbled, but your creativity remains whole.",
			"A quiet mishap; sprinkle sugar on the prompt and begin anew.",
			"Not every song finds its chord on the first stroke, beloved maker.",
		},
	},
}

// GetWhisper selects a pre-written whisper for the persona and trigger category
func GetWhisper(persona string, trigger string) string {
	persona = strings.ToLower(persona)
	personaMap, ok := whispers[persona]
	if !ok {
		// Fallback to victorian_ghost
		personaMap = whispers["victorian_ghost"]
	}

	list, ok := personaMap[trigger]
	if !ok || len(list) == 0 {
		list = personaMap["general_failure"]
	}

	if len(list) == 0 {
		return "A cold draught passes through the terminal."
	}

	idx := rand.Intn(len(list))
	return list[idx]
}
