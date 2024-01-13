package comfyui

var workflow = `{
  "last_node_id": 27,
  "last_link_id": 54,
  "nodes": [
    {
      "id": 7,
      "type": "CLIPTextEncode",
      "pos": [
        352,
        176
      ],
      "size": {
        "0": 425.27801513671875,
        "1": 180.6060791015625
      },
      "flags": {},
      "order": 5,
      "mode": 0,
      "inputs": [
        {
          "name": "clip",
          "type": "CLIP",
          "link": 39,
          "label": "clip"
        }
      ],
      "outputs": [
        {
          "name": "CONDITIONING",
          "type": "CONDITIONING",
          "links": [
            20
          ],
          "slot_index": 0,
          "label": "CONDITIONING"
        }
      ],
      "properties": {
        "Node name for S&R": "CLIPTextEncode"
      },
      "widgets_values": [
        "text, watermark"
      ]
    },
    {
      "id": 5,
      "type": "EmptyLatentImage",
      "pos": [
        462,
        398
      ],
      "size": {
        "0": 315,
        "1": 106
      },
      "flags": {},
      "order": 0,
      "mode": 0,
      "outputs": [
        {
          "name": "LATENT",
          "type": "LATENT",
          "links": [
            23
          ],
          "slot_index": 0,
          "label": "LATENT"
        }
      ],
      "properties": {
        "Node name for S&R": "EmptyLatentImage"
      },
      "widgets_values": [
        512,
        512,
        1
      ]
    },
    {
      "id": 8,
      "type": "VAEDecode",
      "pos": [
        1183,
        -66
      ],
      "size": {
        "0": 210,
        "1": 46
      },
      "flags": {},
      "order": 7,
      "mode": 0,
      "inputs": [
        {
          "name": "samples",
          "type": "LATENT",
          "link": 28,
          "label": "samples"
        },
        {
          "name": "vae",
          "type": "VAE",
          "link": 40,
          "slot_index": 1,
          "label": "vae"
        }
      ],
      "outputs": [
        {
          "name": "IMAGE",
          "type": "IMAGE",
          "links": [
            53,
            54
          ],
          "slot_index": 0,
          "label": "IMAGE"
        }
      ],
      "properties": {
        "Node name for S&R": "VAEDecode"
      }
    },
    {
      "id": 13,
      "type": "SamplerCustom",
      "pos": [
        800,
        -66
      ],
      "size": {
        "0": 355.20001220703125,
        "1": 230
      },
      "flags": {},
      "order": 6,
      "mode": 0,
      "inputs": [
        {
          "name": "model",
          "type": "MODEL",
          "link": 41,
          "slot_index": 0,
          "label": "model"
        },
        {
          "name": "positive",
          "type": "CONDITIONING",
          "link": 19,
          "slot_index": 1,
          "label": "positive"
        },
        {
          "name": "negative",
          "type": "CONDITIONING",
          "link": 20,
          "label": "negative"
        },
        {
          "name": "sampler",
          "type": "SAMPLER",
          "link": 18,
          "slot_index": 3,
          "label": "sampler"
        },
        {
          "name": "sigmas",
          "type": "SIGMAS",
          "link": 49,
          "slot_index": 4,
          "label": "sigmas"
        },
        {
          "name": "latent_image",
          "type": "LATENT",
          "link": 23,
          "slot_index": 5,
          "label": "latent_image"
        }
      ],
      "outputs": [
        {
          "name": "output",
          "type": "LATENT",
          "links": [
            28
          ],
          "shape": 3,
          "slot_index": 0,
          "label": "output"
        },
        {
          "name": "denoised_output",
          "type": "LATENT",
          "links": null,
          "shape": 3,
          "label": "denoised_output"
        }
      ],
      "properties": {
        "Node name for S&R": "SamplerCustom"
      },
      "widgets_values": [
        true,
        0,
        "fixed",
        1
      ]
    },
    {
      "id": 20,
      "type": "CheckpointLoaderSimple",
      "pos": [
        -17,
        -70
      ],
      "size": {
        "0": 343.69647216796875,
        "1": 98
      },
      "flags": {},
      "order": 1,
      "mode": 0,
      "outputs": [
        {
          "name": "MODEL",
          "type": "MODEL",
          "links": [
            41,
            45
          ],
          "shape": 3,
          "slot_index": 0,
          "label": "MODEL"
        },
        {
          "name": "CLIP",
          "type": "CLIP",
          "links": [
            38,
            39
          ],
          "shape": 3,
          "slot_index": 1,
          "label": "CLIP"
        },
        {
          "name": "VAE",
          "type": "VAE",
          "links": [
            40
          ],
          "shape": 3,
          "slot_index": 2,
          "label": "VAE"
        }
      ],
      "properties": {
        "Node name for S&R": "CheckpointLoaderSimple"
      },
      "widgets_values": [
        "sd_xl_turbo_1.0.safetensors"
      ]
    },
    {
      "id": 14,
      "type": "KSamplerSelect",
      "pos": [
        452,
        -144
      ],
      "size": {
        "0": 315,
        "1": 58
      },
      "flags": {},
      "order": 2,
      "mode": 0,
      "outputs": [
        {
          "name": "SAMPLER",
          "type": "SAMPLER",
          "links": [
            18
          ],
          "shape": 3,
          "label": "SAMPLER"
        }
      ],
      "properties": {
        "Node name for S&R": "KSamplerSelect"
      },
      "widgets_values": [
        "dpmpp_2m"
      ]
    },
    {
      "id": 22,
      "type": "SDTurboScheduler",
      "pos": [
        452,
        -248
      ],
      "size": {
        "0": 315,
        "1": 82
      },
      "flags": {},
      "order": 3,
      "mode": 0,
      "inputs": [
        {
          "name": "model",
          "type": "MODEL",
          "link": 45,
          "slot_index": 0,
          "label": "model"
        }
      ],
      "outputs": [
        {
          "name": "SIGMAS",
          "type": "SIGMAS",
          "links": [
            49
          ],
          "shape": 3,
          "slot_index": 0,
          "label": "SIGMAS"
        }
      ],
      "properties": {
        "Node name for S&R": "SDTurboScheduler"
      },
      "widgets_values": [
        2,
        1
      ]
    },
    {
      "id": 6,
      "type": "CLIPTextEncode",
      "pos": [
        351,
        -45
      ],
      "size": {
        "0": 422.84503173828125,
        "1": 164.31304931640625
      },
      "flags": {},
      "order": 4,
      "mode": 0,
      "inputs": [
        {
          "name": "clip",
          "type": "CLIP",
          "link": 38,
          "slot_index": 0,
          "label": "clip"
        }
      ],
      "outputs": [
        {
          "name": "CONDITIONING",
          "type": "CONDITIONING",
          "links": [
            19
          ],
          "slot_index": 0,
          "label": "CONDITIONING"
        }
      ],
      "properties": {
        "Node name for S&R": "CLIPTextEncode"
      },
      "widgets_values": [
        "a boy"
      ]
    },
    {
      "id": 27,
      "type": "SaveImage",
      "pos": [
        1843,
        -154
      ],
      "size": {
        "0": 466.7873840332031,
        "1": 516.8289794921875
      },
      "flags": {},
      "order": 9,
      "mode": 2,
      "inputs": [
        {
          "name": "images",
          "type": "IMAGE",
          "link": 54,
          "label": "images"
        }
      ],
      "properties": {},
      "widgets_values": [
        "ComfyUI"
      ]
    },
    {
      "id": 25,
      "type": "PreviewImage",
      "pos": [
        1213,
        93
      ],
      "size": {
        "0": 501.69647216796875,
        "1": 541.9198608398438
      },
      "flags": {},
      "order": 8,
      "mode": 0,
      "inputs": [
        {
          "name": "images",
          "type": "IMAGE",
          "link": 53,
          "label": "images"
        }
      ],
      "properties": {
        "Node name for S&R": "PreviewImage"
      }
    }
  ],
  "links": [
    [
      18,
      14,
      0,
      13,
      3,
      "SAMPLER"
    ],
    [
      19,
      6,
      0,
      13,
      1,
      "CONDITIONING"
    ],
    [
      20,
      7,
      0,
      13,
      2,
      "CONDITIONING"
    ],
    [
      23,
      5,
      0,
      13,
      5,
      "LATENT"
    ],
    [
      28,
      13,
      0,
      8,
      0,
      "LATENT"
    ],
    [
      38,
      20,
      1,
      6,
      0,
      "CLIP"
    ],
    [
      39,
      20,
      1,
      7,
      0,
      "CLIP"
    ],
    [
      40,
      20,
      2,
      8,
      1,
      "VAE"
    ],
    [
      41,
      20,
      0,
      13,
      0,
      "MODEL"
    ],
    [
      45,
      20,
      0,
      22,
      0,
      "MODEL"
    ],
    [
      49,
      22,
      0,
      13,
      4,
      "SIGMAS"
    ],
    [
      53,
      8,
      0,
      25,
      0,
      "IMAGE"
    ],
    [
      54,
      8,
      0,
      27,
      0,
      "IMAGE"
    ]
  ],
  "groups": [
    {
      "title": "Unmute (CTRL-M) if you want to save images.",
      "bounding": [
        1808,
        -260,
        536,
        676
      ],
      "color": "#3f789e",
      "font_size": 24
    }
  ],
  "config": {},
  "extra": {},
  "version": 0.4
}`
