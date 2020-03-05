package alexakit

import (
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewAlexaRequestIntent_Created(t *testing.T) {
    requestPayload := `
{
  "version": "1.0",
  "session": {
    "new": true,
    "sessionId": "amzn1.echo-api.session.04c4d645-1f5a-44d1-8e28-a2e67680a5f9",
    "application": {
      "applicationId": "amzn1.ask.skill.ae7be50d-3df0-4243-98ca-221e1397fc99"
    },
    "user": {
      "userId": "amzn1.ask.account.AFSN5VT46GYZB4FCHIQCD6WGZ2XYZVS6QE7JUW2A72TQKNPRDB3FAA2DVM2TW3W5EABD5UELAHNBHNQLZ33IA63BMFV6RYQAKKVDPM3USXXRFIKCS35XKXGHGXPTX5DLJQ54ZVZGSPUYID7V7AJ6ZI5LOUE24NRXJF2VBYDCOVSZBQNGG5NHGEIDEPT7UT43FQ7QRUMT37IX2KI"
    }
  },
  "context": {
    "Display": {},
    "System": {
      "application": {
        "applicationId": "amzn1.ask.skill.ae7be50d-3df0-4243-98ca-221e1397fc99"
      },
      "user": {
        "userId": "amzn1.ask.account.AFSN5VT46GYZB4FCHIQCD6WGZ2XYZVS6QE7JUW2A72TQKNPRDB3FAA2DVM2TW3W5EABD5UELAHNBHNQLZ33IA63BMFV6RYQAKKVDPM3USXXRFIKCS35XKXGHGXPTX5DLJQ54ZVZGSPUYID7V7AJ6ZI5LOUE24NRXJF2VBYDCOVSZBQNGG5NHGEIDEPT7UT43FQ7QRUMT37IX2KI"
      },
      "device": {
        "deviceId": "amzn1.ask.device.AERME37LKSLWHXI3UVRSYU5F7K74SFQUQNJU7QBG2XKU7A3MZL35RAEQ7K3UZUZKNP6EHMVO4URAZDK7IFLUG7U2NZNBQJS3IECDK5V2NFVJ27PX4JTXQT4LNOVN2AJE3VAEGECH6XH4E7PWHYDP7DLUAA6NVP6OHLLOTAGWOIKSVAD4WRJDY",
        "supportedInterfaces": {
          "Display": {
            "templateVersion": "1.0",
            "markupVersion": "1.0"
          }
        }
      },
      "apiEndpoint": "https://api.eu.amazonalexa.com",
      "apiAccessToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjEifQ.eyJhdWQiOiJodHRwczovL2FwaS5hbWF6b25hbGV4YS5jb20iLCJpc3MiOiJBbGV4YVNraWxsS2l0Iiwic3ViIjoiYW16bjEuYXNrLnNraWxsLmFlN2JlNTBkLTNkZjAtNDI0My05OGNhLTIyMWUxMzk3ZmM5OSIsImV4cCI6MTU4MzQwMTgyMiwiaWF0IjoxNTgzNDAxNTIyLCJuYmYiOjE1ODM0MDE1MjIsInByaXZhdGVDbGFpbXMiOnsiY29udGV4dCI6IkFBQUFBQUFBQVFEUk44YjJ3WDNWMUtSZ2g3dHZUZGoyS2dFQUFBQUFBQUE1a2EyM2RscU9xL0hxV3BlaE5XajI4KzI1RkxLaWh3TWlsQjdBSExLazd4YXBEbHh1Sk5jcHh6eTlOaFdTU25YWTkyWkFiWWZZK2lyYnp3WnUzcjJwS0wrQW52MHJKTll4dHFaSHVaUkFlREpkQUFVK1p6NVlZZ3ZDVnhvN3BTMnhvNTJuS2d0VlBBZ29BTU00Z2NkbFFXUzJ1SEtTOGloQzYydE1qdU5QZnpHS3MzR0NsUVZ5Nm84MjNDR25PY3F1QndJM1hJNWlPN3JDV1B2ZUsrL1dKVFg4MHJCenJ1MXkwZ0pjTmtEbFBTWWp6R0lablNNdnZHRFhKL1VUZzB4TUJsU3pPM29tSmhNYTBsSkFUTlJXSDRqVlV3dXZNWGt3VGJvK1ArcTNSMWJVUTYwNHlaUGFaK0l1NTB3U1FMQm9LTkFWaThBWXhYc2ZOWEhEdktCL3hkWWRvS3NFYlVoWEtDQ2taUzJ0UjR1elJmMkRMOG9UVWg0YU1WWWhjaWtzOTBUNDBCYkk2WHJYIiwiY29uc2VudFRva2VuIjpudWxsLCJkZXZpY2VJZCI6ImFtem4xLmFzay5kZXZpY2UuQUVSTUUzN0xLU0xXSFhJM1VWUlNZVTVGN0s3NFNGUVVRTkpVN1FCRzJYS1U3QTNNWkwzNVJBRVE3SzNVWlVaS05QNkVITVZPNFVSQVpESzdJRkxVRzdVMk5aTkJRSlMzSUVDREs1VjJORlZKMjdQWDRKVFhRVDRMTk9WTjJBSkUzVkFFR0VDSDZYSDRFN1BXSFlEUDdETFVBQTZOVlA2T0hMTE9UQUdXT0lLU1ZBRDRXUkpEWSIsInVzZXJJZCI6ImFtem4xLmFzay5hY2NvdW50LkFGU041VlQ0NkdZWkI0RkNISVFDRDZXR1oyWFlaVlM2UUU3SlVXMkE3MlRRS05QUkRCM0ZBQTJEVk0yVFczVzVFQUJENVVFTEFITkJITlFMWjMzSUE2M0JNRlY2UllRQUtLVkRQTTNVU1hYUkZJS0NTMzVYS1hHSEdYUFRYNURMSlE1NFpWWkdTUFVZSUQ3VjdBSjZaSTVMT1VFMjROUlhKRjJWQllEQ09WU1pCUU5HRzVOSEdFSURFUFQ3VVQ0M0ZRN1FSVU1UMzdJWDJLSSJ9fQ.luyt_rOYfqcK_fqCZ3rJZVrrzRhtzkUtUeIniXJ8t20vonjpvfwX27P2MEdlEhwaEgxVNsabXRvrGpU_oPj1d2HF1zNcEO1Q60i5nOFQR3j0p-cLJoYBgHY5Ufr6hmw961q5kzv36w59ImQa0JKyt1Ul2sCjG_CzAq3xTi7OTgKBh2b0SRDK85JUfte9u1TJ1-Kib2Bi2I27Fy_hqIDcpnB-JobJ9pYzCeiii_Rkms-9BbPhiLTct6-0QRrBc6OFbtTwf5vO61sQNqOfGvVp5W5-pE0qELBoQK-nJx_EWfWPlwfGlTzoLSpuDlja7jWegFtcNm-TLWAjBKOEDHKpRQ"
    },
    "Viewport": {
      "experiences": [
        {
          "arcMinuteWidth": 246,
          "arcMinuteHeight": 144,
          "canRotate": false,
          "canResize": false
        }
      ],
      "shape": "RECTANGLE",
      "pixelWidth": 1024,
      "pixelHeight": 600,
      "dpi": 160,
      "currentPixelWidth": 1024,
      "currentPixelHeight": 600,
      "touch": [
        "SINGLE"
      ],
      "video": {
        "codecs": [
          "H_264_42",
          "H_264_41"
        ]
      }
    },
    "Viewports": [
      {
        "type": "APL",
        "id": "main",
        "shape": "RECTANGLE",
        "dpi": 160,
        "presentationType": "STANDARD",
        "canRotate": false,
        "configuration": {
          "current": {
            "video": {
              "codecs": [
                "H_264_42",
                "H_264_41"
              ]
            },
            "size": {
              "type": "DISCRETE",
              "pixelWidth": 1024,
              "pixelHeight": 600
            }
          }
        }
      }
    ]
  },
  "request": {
    "type": "IntentRequest",
    "requestId": "amzn1.echo-api.request.489985c5-bb02-45c9-8365-fbc4d2ced01d",
    "timestamp": "2020-03-05T09:45:22Z",
    "locale": "en-US",
    "intent": {
      "name": "TurnOnIntent",
      "confirmationStatus": "NONE",
      "slots": {
        "action": {
          "name": "action",
          "value": "turn on",
          "resolutions": {
            "resolutionsPerAuthority": [
              {
                "authority": "amzn1.er-authority.echo-sdk.amzn1.ask.skill.ae7be50d-3df0-4243-98ca-221e1397fc99.actionType",
                "status": {
                  "code": "ER_SUCCESS_MATCH"
                },
                "values": [
                  {
                    "value": {
                      "name": "on",
                      "id": "ed2b5c0139cec8ad2873829dc1117d50"
                    }
                  },
                  {
                    "value": {
                      "name": "off",
                      "id": "3262d48df5d75e3452f0f16b313b7808"
                    }
                  }
                ]
              }
            ]
          },
          "confirmationStatus": "NONE",
          "source": "USER"
        },
        "item": {
          "name": "item",
          "value": "light",
          "resolutions": {
            "resolutionsPerAuthority": [
              {
                "authority": "amzn1.er-authority.echo-sdk.amzn1.ask.skill.ae7be50d-3df0-4243-98ca-221e1397fc99.ItemType",
                "status": {
                  "code": "ER_SUCCESS_MATCH"
                },
                "values": [
                  {
                    "value": {
                      "name": "light",
                      "id": "2ac43aa43bf473f9a9c09b4b608619d3"
                    }
                  }
                ]
              }
            ]
          },
          "confirmationStatus": "NONE",
          "source": "USER"
        }
      }
    }
  }
}
`
    r := httptest.NewRequest("POST", "/run/intent",strings.NewReader(requestPayload))

    expectedAlexaRequest := AlexaRequest{
        Version: "1.0",
        Request: Request{
            Type:      "IntentRequest",
            RequestID: "amzn1.echo-api.request.489985c5-bb02-45c9-8365-fbc4d2ced01d",
            TimeStamp: "2020-03-05T09:45:22Z",
            Locale:    "en-US",
            Intent:    Intent{
                Name:               "TurnOnIntent",
                ConfirmationStatus: "NONE",
                Slots:      map[string]Slot{
                    "action": {
                        Name:               "action",
                        Value:              "turn on",
                        ConfirmationStatus: "NONE",
                        Source:             "USER",
                    },
                    "item": {
                        Name:               "item",
                        Value:              "light",
                        ConfirmationStatus: "NONE",
                        Source:             "USER",
                    },
                },
            },
        },
    }

    alexaRequest, err := NewAlexaRequestIntent(r)
    assert.Nil(t, err)
    assert.ObjectsAreEqual(expectedAlexaRequest, alexaRequest)
}
