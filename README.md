# Protoschedule
a weekly schedule prototyper package (written in Go)

With this software you can define a prototypical weekly schedule in JSON and then determine if a particular time falls within. Intervals within the schedule (shifts) can be labeled.

This package has no external dependencies.

Example schedule definition:

```json
{
	"description": "Normal Schedule",
	"schedule": {
	     "mon": 
	        [
	            { "label": "shift1", "start": "0800", "duration": "4h" },
	            { "label": "shift1", "start": "1245", "duration": "4h15m" },
	            { "label": "shift2", "start": "1700", "duration": "4h" },
	            { "label": "shift2", "start": "2330", "duration": "4h15m" }
	        ],
	    "tue":
	        [
	            { "label": "shift1", "start": "0800", "duration": "4h" },
	            { "label": "shift1", "start": "1245", "duration": "4h15m" },
	            { "label": "shift2", "start": "1700", "duration": "4h" },
	            { "label": "shift2", "start": "2330", "duration": "4h15m" }
	        ],
	    "wed":
	        [
	            { "label": "shift1", "start": "0800", "duration": "4h" },
	            { "label": "shift1", "start": "1245", "duration": "4h15m" },
	            { "label": "shift2", "start": "1700", "duration": "4h" },
	            { "label": "shift2", "start": "2330", "duration": "4h15m" }
	        ],
	    "thu":
	        [
	            { "label": "shift1", "start": "0800", "duration": "4h" },
	            { "label": "shift1", "start": "1245", "duration": "4h15m" },
	            { "label": "shift2", "start": "1700", "duration": "4h" },
	            { "label": "shift2", "start": "2330", "duration": "4h15m" }
	        ],
	    "fri":
	        [
	            { "label": "shift1", "start": "0800", "duration": "4h" },
	            { "label": "shift1", "start": "1245", "duration": "4h15m" },
	            { "label": "shift2", "start": "1700", "duration": "4h" },
	            { "label": "shift2", "start": "2330", "duration": "4h15m" }
	        ],
	    "sat":
	        [
	            { "label": "shift1", "start": "0900", "duration": "4h" },
	            { "label": "shift2", "start": "1400", "duration": "4h" }
	        ],
	    "sun":
	      []
	  }	    
}
```

See the test fixtures for more examples and protoschedule_test.go for example usage. ```go test -v``` results in some interesting output.

One interesting use of the protoschedule is as a replacement to crontab-style libraries. For instance, the following code could be implemented to invoke a procedure only during the defined schedule.

```golang
	sd, _ := New(jsonString)
	ticker := time.NewTicker(time.Minute)
	go func() {
		for t := range ticker.C {
			if sd.Within(t) {
				SomeFuncThatShouldOnlyBeRunDuringTheDefinedSchedule()
			}
		}
	}()
``` 
